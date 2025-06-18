import { create } from 'zustand';
import { persist } from 'zustand/middleware';

const usePrivateKeyStore = create(
    persist(
        (set) => ({
            privateKey: null,
            keyFile: null,
            setPrivateKey: (key) => set({ privateKey: key }),
            setKeyFile: (file) => set({ keyFile: file }),
            clearAll: () => set({ privateKey: null, keyFile: null }),
            hasKey: () => {
                const state = usePrivateKeyStore.getState();
                return !!state.privateKey || !!state.keyFile;
            }
        })
    )
);

export default usePrivateKeyStore;