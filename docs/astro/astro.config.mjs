import {defineConfig} from 'astro/config';
import starlight from '@astrojs/starlight';
import starlightLinksValidator from 'starlight-links-validator'
import yaml from '@rollup/plugin-yaml';
import sitemap from '@astrojs/sitemap';

// https://astro.build/config
export default defineConfig({
    site: 'https://zfs-backup-blob.top',
    vite: {
        plugins: [yaml()]
    },
    integrations: [
        sitemap(),
        starlight({
            plugins: [starlightLinksValidator()],
            head: [
                {
                    tag: 'link',
                    attrs: {
                        rel: 'icon',
                        href: '/favicon.svg',
                        sizes: '32x32',
                    },
                },
            ],
            title: 'ZFS Backup Blob',
            social: {
                github: 'https://github.com/backup-blob/zfs-backup-blob',
            },
            sidebar: [
                {
                    label: 'Getting started',
                    autogenerate: {directory: 'getting-started'},
                },
                {label: 'Examples', link: '/examples/'},
                {label: 'Architecture', link: '/architecture/'},
                {
                    label: 'Configuration',
                    autogenerate: {directory: 'configuration'},
                    collapsed: true
                },
                {
                    label: 'CLI Commands',
                    autogenerate: {directory: 'cli'},
                    badge: {text: 'Beta', variant: 'caution'},
                    collapsed: true
                },
                {label: 'FAQ', link: '/faq/'},
                {label: 'Feedback', link: '/feedback/'},
            ],
        }),
    ],
});
