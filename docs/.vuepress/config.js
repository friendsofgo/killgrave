import {defaultTheme} from '@vuepress/theme-default'
import {searchPlugin} from '@vuepress/plugin-search'

export default {
	base: '/killgrave/',
	title: 'Killgrave',
	description: 'The simplest way to generate your mock servers',
	head: [['link', {rel: 'icon', href: 'img/killgrave.png'}]],
	plugins: [
		searchPlugin({}),
	],
	theme: defaultTheme({
		logo: 'img/killgrave.png',
		repo: 'friendsofgo/killgrave',
		docsBranch: 'main',
		docsDir: 'docs',
		editLinkPattern: ':repo/edit/:branch/:path',
		themePlugins: {
			backToTop: true,
		},
		navbar: [
			{
				text: 'Get started',
				link: '/guide/',
			},
			{
				text: 'CLI',
				link: '/cli/',
			},
			{
				text: 'Config Reference',
				link: '/config/'
			},
		],
		sidebar: {
			'/guide/': [
				{
					text: 'Introduction',
					link: '/guide/',
					children: [
						'/guide/concepts.md',
						'/guide/installation.md',
					]
				},
				{
					text: 'Using Killgrave',
					link: '/guide/getting-started.md',
					children: [
						'/guide/getting-started.md',
						'/guide/your-first-imposter.md',
						'/guide/advanced.md',
					]
				},
				{
					text: 'How to...?',
					link: '/guide/ht-regex.md',
					children: [
						'/guide/ht-regex.md',
						'/guide/ht-json.md',
						'/guide/ht-delays.md',
						'/guide/ht-dynamic.md',
					]
				}
			],
		},
	}),
}