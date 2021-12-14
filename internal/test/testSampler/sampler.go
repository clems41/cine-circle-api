package testSampler

import (
	"gorm.io/gorm"
	"testing"
)

type Sampler struct {
	t  *testing.T
	DB *gorm.DB
}

// New instancie un Sampler, qui permet de générer de la donnée en base (utile pour les tests unitaires).
// Le but d'utiliser le sampler est d'éviter de faire appel à d'autres domaines que celui testé pour ajouter de la donnée en base.
// Chaque domaine doit être indépendant, et c'est d'autant plus valable pour les tests.
// Example : on veut tester la route qui met à jour un utilisateur, on va d'abord faire appel au sampler pour créer un utilisateur en base plutôt que de faire appel à la route qui créée un utilisateur.
// Ainsi, on s'assure que le test va tester uniquement la route de mise à jour et rien d'autre.
// Il sera cependant possible de faire des tests d'intégration qui appelleront plusieurs routes, mais ce type de test ne doit pas remplacer les tests unitaires.
func New(t *testing.T, DB *gorm.DB, populateDatabase bool) (sampler *Sampler) {

	sampler = &Sampler{
		t:  t,
		DB: DB,
	}

	if populateDatabase {
		sampler.populateDatabase()
	}

	return
}

// populateDatabase si à true, des données vont être insérées avant le lancement du test.
// Cela permet d'avoir une base de donnée non vierge, ce qui se rapproche plus des conditions de production.
func (sampler *Sampler) populateDatabase() {
}
