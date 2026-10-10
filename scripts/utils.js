import path from 'path';
import fs from 'fs/promises';
import child from 'child_process';

export default {
    async findDirFromItems(items) {
        let current = path.resolve(process.cwd());

        while (true) {
            const matches = await Promise.all(items.map(async (item) => {
                return fs.access(path.join(current, item))
                    .then(() => true)
                    .catch(() => false)
            }));

            if (matches.every(Boolean)) {
                return current;
            }

            const parent = path.dirname(current);
            if (parent === current) return null;

            current = parent;
        }
    },

    async runCommand(program, args, options) {
        return new Promise((resolve, reject) => {
            const spawned = child.spawn(program, args, {
                ...options,
                stdio: 'inherit'
            });

            spawned.on('error', reject);

            spawned.on('close', (code, signal) => {
                if (code === 0) {
                    resolve();
                    return;
                }

                reject(new Error(`Error code ${code}${signal ? `, Signal ${signal}` : ''}`))
            });
        });
    },

    bytesToHexString(bytes) {
        return '0x' + bytes
            .toString(16)
            .toUpperCase()
            .padStart(2, '0');
    },

    getCodePart(line) {
        const index = line.indexOf('//');
        return index === -1 ? line : line.slice(0, index);
    }
};
