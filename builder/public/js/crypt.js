export default {
    nonceLength: 12,

    async generateKeyPair() {
        return await crypto.subtle.generateKey(
            { name: 'ECDH', namedCurve: 'P-256' },
            true,
            [ 'deriveKey', 'deriveBits' ]
        );
    },

    async exportPublicKey(publicKey) {
        const exportedKey = await crypto.subtle.exportKey('raw', publicKey);
        return btoa(String.fromCharCode(...new Uint8Array(exportedKey)));
    },

    async getSharedSecret(privateKey, serverPublicKey) {
        const decodedServerPublicKey = Uint8Array.from(
            atob(serverPublicKey),
            (char) => char.charCodeAt(0)
        );

        const theServerPublicKey = await crypto.subtle.importKey(
            'raw',
            decodedServerPublicKey,
            { name: 'ECDH', namedCurve: 'P-256' },
            false,
            []
        );

        return await crypto.subtle.deriveBits(
            { name: 'ECDH', public: theServerPublicKey },
            privateKey,
            256
        );
    },

    async encryptData (data, key) {
        const keyHash = await this.hashKey(key);
        const encoded = new TextEncoder().encode(data);

        const nonce = new Uint8Array(this.nonceLength);
        crypto.getRandomValues(nonce);

        const gcm = await crypto.subtle.encrypt(
            { name: 'AES-GCM', iv: nonce },
            keyHash,
            encoded
        );

        const result = new Uint8Array(nonce.length + gcm.byteLength);
        result.set(nonce, 0);
        result.set(new Uint8Array(gcm), nonce.length);

        return btoa(String.fromCharCode(...result));
    },

    async decryptData(data, key) {
        const keyHash = await this.hashKey(key);
        const decoded = Uint8Array.from(atob(data), (char) => char.charCodeAt(0));

        if (decoded.length < this.nonceLength) {
            throw new Error('data is too short');
        }

        const nonce = decoded.slice(0, this.nonceLength);
        const ciphered = decoded.slice(this.nonceLength);

        const decrypted = await crypto.subtle.decrypt(
            { name: 'AES-GCM', iv: nonce },
            keyHash,
            ciphered
        );

        return new TextDecoder().decode(decrypted);
    },

    async hashKey(key) {
        let keyBytes;

        if (key instanceof CryptoKey) {
            return key;
        } else if (key instanceof ArrayBuffer || ArrayBuffer.isView(key)) {
            keyBytes = key;
        } else {
            keyBytes = new TextEncoder().encode(key);
        }

        const digest = await crypto.subtle.digest('SHA-256', keyBytes);

        return crypto.subtle.importKey(
            'raw',
            digest,
            { name: 'AES-GCM' },
            false,
            ['encrypt', 'decrypt']
        );
    }
};
