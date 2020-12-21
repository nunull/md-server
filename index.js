#!/usr/bin/env node

const express = require('express')
const fs = require('fs')
const showdown  = require('showdown')
const util = require('util')

const readdir = util.promisify(fs.readdir)
const readFile = util.promisify(fs.readFile)

const port = process.env.PORT || 3000

const app = express()
const converter = new showdown.Converter()

app.get('/', (req, res) => {
  readdir('.').then(dir => {
    res.send(page(files(dir.filter(f => f.indexOf('.md') === f.length-3))))
  })
})

app.get('/:path.md', (req, res) => {
  readFile(`./${req.params.path}.md`, { encoding: 'utf-8' }).then(markdown => {
    const html = converter.makeHtml(markdown)
    res.send(page(html, req.params.path + '.md'))
  })
})

app.use(express.static('.'))

// app.get('/:path', (req, res) => {
//   readFile(`./${req.params.path}`).then(raw => {
//     res.send(raw)
//   }).catch(e => console.log(e))
// })


app.listen(port, () => console.log(`listening on :${port}`))

function page (body, path) {
  return `
    <!DOCTYPE html>
    <html lang="en" dir="ltr">
      <head>
        <meta charset="utf-8">
        <title>md-server</title>
        <style>
          * {
            font-family: sans-serif;
          }
          body {
            max-width: 800px;
            margin: 15px auto;
          }
          h1, h2, h3 {
            font-weight: normal;
          }
          img, video {
            display: block;
            max-width: 800px;
            max-height: 80vh;
            margin: 15px auto;
          }
        </style>
      </head>
      <body>
        <nav><a href="/">index</a>${path ? ` <span>${path}</span>` : ''}</nav>
        ${body}
      </body>
    </html>
  `
}

function files (files) {
  return `
    <ul>
      ${files.map(f => `<li><a href="/${f}">${f}</a></li>`).join('')}
    </ul>
  `
}
