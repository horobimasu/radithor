import crypt from './crypt'
import notif from './notif';
import saving from './saving';
import vars from './vars';

const build = document.querySelector('#build');

export default {
    originTitle: document.title,

    timer: null,

    setTitle(title) {
        if (this.originTitle === 'localhost') {
            return;
        }

        document.title = this.originTitle + ' - ' + title;

        if (this.timer) {
            clearTimeout(this.timer);
        }

        this.timer = setTimeout(() => {
            document.title = this.originTitle;
            this.timer = null;
        }, 6 * 60 * 1000);
    },

    initBuild() {
        build.addEventListener('click', async () => {
            let config = await saving.getConfig();
            let parsedConfig = JSON.parse(config);

            if (parsedConfig.end_to_end_encryption === true) {
                const clientKeys = await crypt.generateKeyPair();

                const fetchedServerPublicKey = await fetch('/api/get-key', {
                    method: 'POST',
                    headers: { 'Content-Type': 'application/json' },
                    body: JSON.stringify({ key: await crypt.exportPublicKey(clientKeys.publicKey) })
                });

                if (!fetchedServerPublicKey.ok) {
                    const text = await fetchedServerPublicKey.text();

                    notif.showNotif('Failed to get server key:\n' + text + '\n\n' + 'See console for full error.', vars.NOTIF_ERROR, 10);
                    console.warn(text);
                    return;
                }

                const serverPublicKey = await fetchedServerPublicKey.text();
                const sharedSecret = await crypt.getSharedSecret(clientKeys.privateKey, serverPublicKey);

                config = await saving.getConfig(sharedSecret);
            }

            const initial = notif.showNotif('Building...', vars.NOTIF_INFO, 900);
            this.setTitle('Building');

            const build = await fetch('/api/gen-build', {
                method: 'POST',
                headers: { 'Content-Type': 'application/json' },
                body: config
            });

            if (!build.ok) {
                notif.removeNotif(initial);
                this.setTitle('Failed');

                const text = await build.text();

                notif.showNotif('Failed to build:\n' + text, vars.NOTIF_ERROR, 30);
                console.warn(text);
                return;
            }

            notif.removeNotif(initial);
            notif.showNotif('Successfully built.', vars.NOTIF_INFO, 90);
            this.setTitle('Success');
        });
    }
};
