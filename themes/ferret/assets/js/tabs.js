(function () {
    function getTabParts(tabset) {
        const list = tabset.querySelector('.docs-tabs-list');
        const tabs = list ? Array.from(list.querySelectorAll('.docs-tabs-tab')) : [];
        const panels = Array.from(tabset.children).filter((child) => child.classList.contains('docs-tabs-panel'));

        return { tabs, panels };
    }

    function resizeEditors(panel) {
        if (!panel) {
            return;
        }

        if (window.ace) {
            panel.querySelectorAll('.ace_editor').forEach((editorEl) => {
                if (editorEl.env && editorEl.env.editor) {
                    editorEl.env.editor.resize(true);
                }
            });
        }

        window.dispatchEvent(new Event('resize'));
    }

    function activateTab(tabset, nextIndex, shouldFocus) {
        const { tabs, panels } = getTabParts(tabset);

        if (!tabs[nextIndex] || !panels[nextIndex]) {
            return;
        }

        tabs.forEach((tab, index) => {
            const selected = index === nextIndex;

            tab.setAttribute('aria-selected', selected ? 'true' : 'false');
            tab.tabIndex = selected ? 0 : -1;
        });

        panels.forEach((panel, index) => {
            panel.hidden = index !== nextIndex;
        });

        if (shouldFocus) {
            tabs[nextIndex].focus();
        }

        window.requestAnimationFrame(() => resizeEditors(panels[nextIndex]));
    }

    function getActiveIndex(tabs) {
        const selectedIndex = tabs.findIndex((tab) => tab.getAttribute('aria-selected') === 'true');

        return selectedIndex >= 0 ? selectedIndex : 0;
    }

    function initTabs(tabset) {
        const { tabs } = getTabParts(tabset);

        if (tabs.length === 0) {
            return;
        }

        tabs.forEach((tab, index) => {
            tab.addEventListener('click', () => {
                activateTab(tabset, index, false);
            });

            tab.addEventListener('keydown', (event) => {
                const activeIndex = getActiveIndex(tabs);
                let nextIndex = activeIndex;

                if (event.key === 'ArrowLeft') {
                    nextIndex = activeIndex === 0 ? tabs.length - 1 : activeIndex - 1;
                } else if (event.key === 'ArrowRight') {
                    nextIndex = activeIndex === tabs.length - 1 ? 0 : activeIndex + 1;
                } else if (event.key === 'Home') {
                    nextIndex = 0;
                } else if (event.key === 'End') {
                    nextIndex = tabs.length - 1;
                } else if (event.key !== 'Enter' && event.key !== ' ') {
                    return;
                }

                event.preventDefault();
                activateTab(tabset, nextIndex, true);
            });
        });

        activateTab(tabset, getActiveIndex(tabs), false);
    }

    document.addEventListener('DOMContentLoaded', () => {
        document.querySelectorAll('.docs-tabs').forEach(initTabs);
    }, false);
})();
