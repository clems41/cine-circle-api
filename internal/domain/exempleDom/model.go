package exempleDom

import (
	"cine-circle-api/pkg/utils/searchUtils"
)

/* Common */

type CommonForm struct {
	// TODO add your custom fields here and their rules (cf. userDom example with validate tag)
}

type CommonView struct {
	Id uint `json:"id"`
	// TODO add your custom fields here
}

/* Create */

type CreateForm struct {
	CommonForm
}

type CreateView struct {
	CommonView
}

/* Update */

type UpdateForm struct {
	ExempleId uint `json:"-"` // Champ récupéré depuis le path parameter de la route
	CommonForm
}

type UpdateView struct {
	CommonView
}

/* Delete */

type DeleteForm struct {
	ExempleId uint `json:"-"` // Champ récupéré depuis le path parameter de la route
}

/* Get */

type GetForm struct {
	ExempleId uint `json:"-"` // Champ récupéré depuis le path parameter de la route
}

type GetView struct {
	CommonView
}

/* Search */

// SearchForm permet de paginer, trier et filtrer les exemples.
// Les tags sont utilisés par la méthode qui récupère les query parameters pour remplir cette structure.
// Pour les filtres par tri, les tags utilisés correspondent aux tags des champs habituellement utilisés par le front.
type SearchForm struct {
	searchUtils.PaginationRequest
	searchUtils.SortingRequest
	// TODO add your keyword fields here (cf. userDom example)
}

type SearchView struct {
	searchUtils.Page
	Exemples []CommonView `json:"exemples"`
}
