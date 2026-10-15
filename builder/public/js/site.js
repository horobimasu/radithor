import build from './build';
import crypt from './crypt';
import saving from './saving';
import vars from './vars';

const layout = 'LIVE_BUILDER_SITE_LAYOUT';

const tabs = document.querySelectorAll('.tab-button');
const contents = document.querySelectorAll('.content');

export default {
    createCheckbox(config) {
        const box = document.createElement('div');
        box.className = 'box';

        const title = document.createElement('div');
        title.className = 'title';

        const id = config.name.toLowerCase().replaceAll(' ', '_');

        const input = document.createElement('input');
        input.type = 'checkbox';
        input.id = config.id || id;
        input.checked = config.value ?? false;

        const label = document.createElement('label');
        label.innerHTML = `<strong>${config.name}</strong>`;

        title.appendChild(input);
        title.appendChild(label);

        const desc = document.createElement('div');
        desc.className = 'desc';
        desc.textContent = config.info;

        box.appendChild(title);
        box.appendChild(desc);

        return box;
    },

    createTextbox(config) {
        const box = document.createElement('div');
        box.className = 'box';

        const title = document.createElement('div');
        title.className = 'title';

        const id = config.name.toLowerCase().replaceAll(' ', '_');

        const label = document.createElement('label');
        label.innerHTML = `<strong>${config.name}</strong>`;

        title.appendChild(label);

        const desc = document.createElement('div');
        desc.className = 'desc';
        desc.textContent = config.info;

        const container = document.createElement('div');
        container.className = 'input-container';

        const input = document.createElement('input');
        input.className = 'input';
        input.id = config.id || id;
        input.type = 'text';
        input.value = config.value ?? '';

        if (config.attributes) {
            for (const [key, val] of Object.entries(config.attributes)) {
                input.setAttribute(key, val);
            }
        }

        container.appendChild(input);

        box.appendChild(title);
        box.appendChild(desc);
        box.appendChild(container);

        return box;
    },

    createDropdown(config) {
        const box = document.createElement('div');
        box.className = 'box';

        const title = document.createElement('div');
        title.className = 'title';

        const id = config.name.toLowerCase().replaceAll(' ', '_');

        const label = document.createElement('label');
        label.innerHTML = `<strong>${config.name}</strong>`;

        title.appendChild(label);

        const desc = document.createElement('div');
        desc.className = 'desc';
        desc.textContent = config.info;

        const container = document.createElement('div');
        container.className = 'input-container';

        const select = document.createElement('select');
        select.className = 'input';
        select.id = config.id || id;

        if (config.items && Array.isArray(config.items)) {
            config.items.forEach((item) => {
                const option = document.createElement('option');
                option.value = item;
                option.textContent = item;
                select.appendChild(option);
            });
        }

        container.appendChild(select);

        box.appendChild(title);
        box.appendChild(desc);
        box.appendChild(container);

        return box;
    },

    createSlider(config) {
        const box = document.createElement('div');
        box.className = 'box';

        const title = document.createElement('div');
        title.className = 'title';

        const id = config.name.toLowerCase().replaceAll(' ', '_');

        const label = document.createElement('label');
        label.innerHTML = `<strong>${config.name}</strong>`;

        const value = document.createElement('input');
        value.className = 'value-input';
        value.type = 'number';
        value.id = config.id || id;
        value.value = config.value ?? config.min ?? 0;

        value.addEventListener('blur', () => {
            const val = parseFloat(value.value);
            if (!isNaN(val)) return;

            const theDefault = config.value ?? config.min ?? 0;
            value.value = theDefault;
            slider.value = theDefault;
        });

        value.addEventListener('input', () => {
            const val = parseFloat(value.value);
            if (isNaN(val)) return;

            value.value = Math.min(val, config.max ?? 100);
            slider.value = value.value;
        });

        title.appendChild(label);
        title.appendChild(value);

        const desc = document.createElement('div');
        desc.className = 'desc';
        desc.textContent = config.info;

        const container = document.createElement('div');
        container.className = 'input-container';

        const slider = document.createElement('input');
        slider.type = 'range';
        slider.className = 'input';
        slider.id = config.id || id;
        slider.min = config.min ?? 0;
        slider.max = config.max ?? 100;
        slider.value = config.value ?? config.min ?? 0;

        slider.addEventListener('input', () => {
            value.value = slider.value;
        });

        container.appendChild(slider);

        box.appendChild(title);
        box.appendChild(desc);
        box.appendChild(container);

        return box;
    },

    createSeparator() {
        const separator = document.createElement('div');
        separator.className = 'separator';

        return separator;
    },

    createSetting(config) {
        switch (config.type) {
            case 'separator':
                return this.createSeparator();

            case 'checkbox':
                return this.createCheckbox(config);

            case 'text':
                return this.createTextbox(config);

            case 'slider':
                return this.createSlider(config);

            case 'list':
                return this.createDropdown(config);
        }
    },

    async initSite() {
        tabs.forEach((tab) => {
            tab.addEventListener('click', () => {
                tabs.forEach((button) => button.classList.remove('active'));
                contents.forEach((content) => content.classList.remove('active'));

                tab.classList.add('active');

                const target = tab.getAttribute('data-tab');
                document.getElementById(`${target}-tab`).classList.add('active');
            });
        });

        for (const [tab, item] of Object.entries(layout)) {
            const parent = document.getElementById(`${tab}-tab`);

            item.forEach((config) => {
                const box = this.createSetting(config);
                parent.appendChild(box);
            });
        }

        document.querySelectorAll('input').forEach((input) => {
            input.setAttribute('spellcheck', 'false');
            input.setAttribute('autocomplete', 'off');
        });

        saving.loadConfig(vars.e2eePassword);

        const inputs = document.querySelectorAll('input');
        const selects = document.querySelectorAll('select');

        [...inputs, ...selects].forEach((input) => {
            input.addEventListener('input', async (e) => {
                if (!vars.savingDisabled) {
                    localStorage.setItem(vars.CFG_SAVE_NAME, await saving.getConfig(vars.e2eePassword));
                }

                const target = e.target;
                const value = (() => {
                    if (target.type === 'checkbox') {
                        return target.checked;
                    } else {
                        return target.value.trim();
                    }
                })();

                if (target.id === 'end_to_end_encryption') {
                    if (value) {
                        if (localStorage.getItem(vars.E2EE_SAVE_NAME)) {
                            alert('End to end encryption is already active.');
                            return;
                        }

                        const password = prompt('Enter password to turn on end to end encryption');
                        if (password === '' || !password) {
                            notif('Password is empty or cancelled.');

                            target.checked = false;
                            localStorage.setItem(vars.CFG_SAVE_NAME, await saving.getConfig());
                            return;
                        }

                        vars.e2eePassword = password;

                        localStorage.setItem(vars.CFG_SAVE_NAME, await saving.getConfig(vars.e2eePassword));

                        // this is for when the page reloads and the saved config has to be decrypted
                        // this will verify is the password is correct
                        localStorage.setItem(
                            vars.E2EE_SAVE_NAME,
                            await crypt.encryptData(vars.E2EE_VERIFICATION_STRING, vars.e2eePassword)
                        );
                    } else {
                        localStorage.setItem(vars.CFG_SAVE_NAME, await saving.getConfig());
                        localStorage.setItem(vars.TAG_SAVE_NAME, document.querySelector('.build-tag').textContent);
                        localStorage.removeItem(vars.E2EE_SAVE_NAME);

                        vars.e2eePassword = null;
                    }
                } else if (target.id === 'save_config') {
                    target.checked = !vars.savingDisabled;

                    vars.savingDisabled = !value;

                    if (value) {
                        localStorage.setItem(vars.CFG_SAVE_NAME, await saving.getConfig(vars.e2eePassword));
                    } else {
                        localStorage.setItem(vars.CFG_SAVE_NAME, false);
                    }
                } else if (target.id === 'obfuscate_site_title') {
                    document.title = target.checked ? 'localhost' : 'Radithor builder';
                    build.originTitle = target.checked ? 'localhost' : 'Radithor builder';
                }
            });
        });
    },

    async preSiteInit() {
        vars.savingDisabled = localStorage.getItem(vars.CFG_SAVE_NAME) === 'false';

        const e2eeVerification = localStorage.getItem(vars.E2EE_SAVE_NAME);
        if (e2eeVerification) {
            const password = prompt('Enter password to decrypt data');
            if (password === '') {
                alert('Password is empty.');
                location.reload();
            }

            vars.e2eePassword = password;

            (async () => {
                await crypt.decryptData(e2eeVerification, vars.e2eePassword)
                    .catch(() => {
                        alert('Password is invalid.');
                        location.reload();
                    });
            })();
        }

        (async () => {
            const fetched = await fetch('/api/get-version');
            const versions = await fetched.text();

            const split = versions.split(' ')

            window.serverVer = split[0];
            window.siteVer = split[1];
        })();
    }
};
