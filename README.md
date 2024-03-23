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

By default, omnivore-exporter uses [obelisk](https://github.com/go-shiori/obelisk) to export the articles to a single HTML file on disk, without additional dependencies.

It can optionally use [monolith](https://github.com/Y2Z/monolith) instead of obelisk, if available in PATH. The the fidelity of the exports with monolith is generally better, in its current version.

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

## User service

You can use the [omnivore-exporter.service](/extra/omnivore-exporter.service) and [omnivore-exporter.timer](/extra/omnivore-exporter.timer) to run the exporter as a user service. Copy the files to `~/.config/systemd/user/` and enable it with:

```bash
systemctl --user enable --now omnivore-export.timer
```

Defaults to export every hour, you can change the timer to your liking.

## Credits

- [go-shiori/obelisk](https://github.com/go-shiori/obelisk) - Awesome tool and Go library to export URL contents to disk.
- [Y2Z/monolith](https://github.com/Y2Z/monolith) - Another awesome tool to export URL contents to disk, that inspired Obelisk.
