const { description } = require('../../package')

module.exports = {
  /**
   * Ref：https://v1.vuepress.vuejs.org/config/#title
   */
  title: 'Killgrave',
  /**
   * Ref：https://v1.vuepress.vuejs.org/config/#description
   */
  description: description,

  /**
   * Extra tags to be injected to the page HTML `<head>`
   *
   * ref：https://v1.vuepress.vuejs.org/config/#head
   */
  head: [
    ['meta', { name: 'theme-color', content: '#db512a' }],
    ['meta', { name: 'apple-mobile-web-app-capable', content: 'yes' }],
    ['meta', { name: 'apple-mobile-web-app-status-bar-style', content: 'black' }]
  ],

  /**
   * Theme configuration, here is the default theme configuration for VuePress.
   *
   * ref：https://v1.vuepress.vuejs.org/theme/default-theme-config.html
   */
  themeConfig: {
    smoothScroll: true,
    repo: 'friendsofgo/killgrave',
    docsDir: 'docs',
    editLinks: true,
    prevLinks: true,
    nextLinks: true,
    activeHeaderLinks: true,
    displayAllHeaders: true,
    editLinkText: 'Edit this page on GitHub',
    lastUpdated: 'Last Updated',
    nav: [
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
          title: 'Introduction',
          collapsable: false,
          children: [
            '',
            'concepts',
            'installation',
          ]
        },
        {
          title: 'Using Killgrave',
          collapsable: false,
          children: [
            'getting-started',
            'your-first-imposter',
            'advanced',
          ]
        },
        {
          title: 'How to...?',
          collapsable: false,
          children: [
            'ht-regex',
            'ht-json',
            'ht-delays',
            'ht-dynamic',
          ]
        },
        {
          title: 'Interactive mode',
          collapsable: false,
          children: [
            'debug-intro.md'
          ]
        }
      ],
    }
  },

  /**
   * Apply plugins，ref：https://v1.vuepress.vuejs.org/zh/plugin/
   */
  plugins: [
    '@vuepress/plugin-back-to-top',
    '@vuepress/plugin-last-updated',
    '@vuepress/google-analytics',
    {
      'ga': '' // UA-00000000-0
    },
    'vuepress-plugin-clean-urls',
    {
      normalSuffix: '/',
      indexSuffix: '/',
      notFoundPath: '/404.html',
    }
  ]
}
