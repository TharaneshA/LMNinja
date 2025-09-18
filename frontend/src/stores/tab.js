import { defineStore } from 'pinia';
import { findIndex, get } from 'lodash';
import useBrowserStore from 'stores/browser.js';

class TabItem {
    constructor({ name, db = 0 }) {
        this.name = name; // Connection Name, e.g., "Personal OpenAI"
        this.title = name;
        this.db = db;
        this.subTab = 'chat'; // Default view
    }
}

const useTabStore = defineStore('tab', {
    state: () => ({
        nav: 'server', // 'server' (connection list) or 'browser' (active connections)
        tabList: [],
        activatedIndex: -1,
    }),
    getters: {
        tabs: (state) => state.tabList,
        currentTab: (state) => state.tabList[state.activatedIndex],
        currentTabName: (state) => get(state.tabList[state.activatedIndex], 'name'),
    },
    actions: {
        _setActivatedIndex(idx, switchNav = false) {
            this.activatedIndex = idx;
            if (switchNav) {
                this.nav = idx >= 0 ? 'browser' : 'server';
            }
        },

        switchTab(tabIndex) {
            if (this.activatedIndex !== tabIndex) {
                this._setActivatedIndex(tabIndex);
            }
        },

        upsertTab({ server, db = 0, forceSwitch = false }) {
            let tabIndex = findIndex(this.tabList, { name: server });
            if (tabIndex === -1) {
                const tabItem = new TabItem({ name: server, db });
                this.tabList.push(tabItem);
                tabIndex = this.tabList.length - 1;
            }
            if (forceSwitch) {
                this._setActivatedIndex(tabIndex, true);
            }
        },

        closeTab(tabName) {
            const browserStore = useBrowserStore();
            browserStore.closeConnection(tabName);
        },

        removeTabByName(tabName) {
            const idx = findIndex(this.tabs, { name: tabName });
            if (idx !== -1) {
                this.tabList.splice(idx, 1);
                if (this.activatedIndex >= idx && this.activatedIndex > 0) {
                    this.activatedIndex -= 1;
                } else if (this.tabList.length === 0) {
                    this._setActivatedIndex(-1, true); // No tabs left, switch to server view
                }
                // If the last tab was closed, activatedIndex might now be out of bounds if it was 0
                if(this.activatedIndex >= this.tabList.length) {
                    this.activatedIndex = this.tabList.length - 1;
                }
            }
        },
    },
});

export default useTabStore;