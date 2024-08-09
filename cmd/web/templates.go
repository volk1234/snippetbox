package main

import "snippetbox.sandx.link/internal/models"

type templateData struct {
	Snippet *models.Snippet
	Snippets []*models.Snippet
}
