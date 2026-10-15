import vars from './vars';
import crypt from './crypt';

export default {
    async getConfig(encryptionKey) {
        const config = {
            tag: localStorage.getItem(vars.TAG_SAVE_NAME)
        };

        const inputs = document.querySelectorAll('input');
        const selects = document.querySelectorAll('select');

        [...inputs, ...selects].forEach((input) => {
            if (input.type === 'checkbox') {
                config[input.id] = input.checked;
            } else {
                config[input.id] = input.value.trim();
            }
        });

        if (encryptionKey) {
            const encryptedConfig = Object.entries(config).map(async ([name, value]) => {
                const newValue = await crypt.encryptData(JSON.stringify(value), encryptionKey);
                return [name, newValue];
            });

            const newConfig = Object.fromEntries(
                await Promise.all(encryptedConfig)
            );

            return JSON.stringify(newConfig);
        }

        return JSON.stringify(config);
    },

    async loadConfig(encryptionKey) {
        const saved = localStorage.getItem(vars.CFG_SAVE_NAME);
        if (!saved) return;

        let config = JSON.parse(saved);

        if (encryptionKey) {
            const decryptedConfig = Object.entries(config).map(async ([name, value]) => {
                const decrypted = await crypt.decryptData(String(value), encryptionKey);

                try {
                    return [name, JSON.parse(decrypted)];
                } catch {
                    return [name, decrypted];
                }
            });

            config = Object.fromEntries(
                await Promise.all(decryptedConfig)
            );
        }

        const inputs = document.querySelectorAll('input');
        const selects = document.querySelectorAll('select');

        [...inputs, ...selects].forEach((input) => {
            if (input.type === 'checkbox') {
                input.checked = config[input.id] === true;
            } else {
                input.value = config[input.id] ?? '';
            }
        });
    }
};
