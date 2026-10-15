import crypt from './crypt';
import vars from './vars';

const tag = document.querySelector('.build-tag');
const saved = localStorage.getItem(vars.TAG_SAVE_NAME) || '';

export default {
    async initTag() {
        tag.addEventListener('dblclick', (e) => {
            if (!e.ctrlKey) return;

            tag.contentEditable = true;
            tag.focus();
        });

        tag.addEventListener('blur', async () => {
            tag.contentEditable = false;
            const text = tag.textContent.trim();

            if (text) {
                if (vars.e2eePassword) {
                    const encrypted = await crypt.encryptData(text, vars.e2eePassword)
                    localStorage.setItem(vars.TAG_SAVE_NAME, encrypted);
                } else {
                    localStorage.setItem(vars.TAG_SAVE_NAME, text);
                }
            } else {
                tag.textContent = '';
                localStorage.removeItem(vars.TAG_SAVE_NAME);
            }
        });

        tag.addEventListener('keydown', (e) => {
            if (e.key !== 'Enter') return;

            e.preventDefault();
            tag.blur();
        });

        if (saved) {
            if (vars.e2eePassword) {
                const decrypted = await crypt.decryptData(saved, vars.e2eePassword)
                tag.textContent = decrypted;
            } else {
                tag.textContent = saved;
            }
        }

        if (tag.textContent !== 'default') {
            return;
        }

        tag.style.visibility = 'hidden';

        document.addEventListener('keydown', (e) => {
            if (!e.ctrlKey) return;

            tag.style.visibility = 'visible';
        });
    }
};
