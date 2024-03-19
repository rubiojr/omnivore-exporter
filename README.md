# Omnivore Export

Export [omnivore.app](https://omnivore.app) articles to markdown.

## Installation

```bash
go install github.com/rubiojr/omnivore-export@latest
```

## Usage

[Get your API key](https://docs.omnivore.app/integrations/api.html#getting-an-api-token) and export it:
```bash
export OMNIVORE_API_KEY=your-api-key
```

```bash
omnivore-export --output-dir ~/Documents/omnivore-exports
```

omnivore-exporter can also use [monolith](https://github.com/Y2Z/monolith) to export the articles to a single HTML file, if available in the PATH (recommended, the fidelity of the exports is better). 

```bash
omnivore-export --output-dir ~/Documents/omnivore-exports --use-monolith
```

By default, all articles except the ones labeled with `RSS`, `Newsletter` or `omnivore-exporter-skip` are exported. Select the labeled articles to export with the `--labels` flag:

```bash
omnivore-export --output-dir ~/Documents/omnivore-exports --labels label-to-export --labels another-label
```

You can also use `--skip-labels` to exclude articles with specific labels and export the rest:

```bash
omnivore-export --output-dir ~/Documents/omnivore-exports --skip-labels label-to-skip
```
