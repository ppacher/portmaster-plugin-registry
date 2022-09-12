import { defaultTheme, defineUserConfig } from 'vuepress';

export default defineUserConfig({
  lang: 'en-US',
  title: 'PECS',
  description: 'An easy way to manage your third-party Portmaster Plugins',
  theme: defaultTheme({
    repo: 'ppacher/portmaster-plugin-registry',
    docsRepo: 'https://github.com/ppacher/portmaster-plugin-registry',
    docsBranch: 'main',
    docsDir: 'docs',
    editLinkPattern: ':repo/-/edit/:branch/:path',
    navbar: [
      {
        text: 'Get Started',
        link: '/guide/getting-started',
      },
      {
        text: 'Publish',
        link: 'publish',
      },
      {
        text: 'Developer docs',
        link: 'devdocs',
      },
      {
        text: 'Report an issue',
        link: 'https://github.com/ppacher/portmaster-plugin-registry/issues'
      }
    ]
  }),
})