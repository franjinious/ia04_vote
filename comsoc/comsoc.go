//
// Created by hongzhe on 22-10-3.
// Implémentation des fonctions de base du serveur web de vote.
//

package comsoc

import "errors"

type Alternative int // Candidat
type Profile [][]Alternative // Profile du tous les candidats, Profile[I][J] représente le Jième candidat préféré du Ième votant
type Count map[Alternative]int // Note pour chaque candidat

/**
 * rank
 * @Description: trouver l'indice d'un candidat dans un préférence
 * @param alt: un candidat
 * @param prefs: préférence du un votant
 * @return int: indice de ce candiat dans le préférence
 */
func rank(alt Alternative, prefs []Alternative) int {
	for i := 0; i < len(prefs); i++ {
		if prefs[i] == alt {
			return i
		}
	}
	return -1
}

/**
 * isPref
 * @Description: comparer l'ordre de deux candidats dans une préférence
 * @param alt1: candidat un
 * @param alt2: candidat deux
 * @param prefs: préférence du un votant
 * @return bool: 1 signifie alt1 gagne et 0 signifie alt2 gagne
 */
func isPref(alt1, alt2 Alternative, prefs []Alternative) bool {
	return rank(alt1,prefs) > rank(alt2,prefs)
}

/**
 * maxCount
 * @Description: Comparer les résultats de tous les candidats pour trouver le meilleur
 * @param count: note pour chaque candidat
 * @return bestAlts: un slice du gagant
 */
func maxCount(count Count) (bestAlts []Alternative){
	max_value := -1

	// Trouver le plus grand nombre de votes
	for _,j := range count {
		if j > max_value {
			max_value = j
		}
	}

	// Ajouter des candidats égaux à la valeur maximale à la valeur de retour
	for i,j := range count {
		if j == max_value {
			bestAlts = append(bestAlts, i)
		}
	}

	return bestAlts
}

/**
 * checkProfile
 * @Description: vérifier les règles de préférence
 * @param prefs: profile du tous les candidats
 * @return error: la renvoie s'il y a une erreur, nil s'il n'y a pas d'erreur
 */
func checkProfile(prefs Profile) error {

	// vérifier si il est vide
	if len(prefs) == 0 {
		return errors.New("Le profile list est NULL!")
	}

	// Vérifier si toutes les longueurs sont cohérentes
	length := len(prefs[0])

	for i := 0; i < len(prefs); i++ {
		if len(prefs[i]) != length {
			return errors.New("Profils ne sont pas tous complets!")
		}
	}

	// un Set pour tous les candidats
	set := make(map[Alternative]bool)
	for _,j := range prefs[0] {
		if set[j] == true {
			return errors.New("Pas le seul candidat dans une préférence!")
		}
		set[j] = true
	}

	// Pour chaque votant, déterminez s'il y a des membres en double
	for i := 1; i < len(prefs); i++ {
		set_temp := set
		for _,j := range prefs[i] {
			if set_temp[j] == false {
				return errors.New("Pas le seul candidat dans une préférence!")
			}
			set_temp[j] = false
		}
	}

	return nil
}

/**
 * checkProfileAlternative
 * @Description: vérifier les règles de préférence selon le slice du cadidats
 * @param prefs: profile du tous les candidats
 * @param alts: tous les candidats
 * @return error: la renvoie s'il y a une erreur, nil s'il n'y a pas d'erreur
 */
func checkProfileAlternative(prefs Profile, alts []Alternative) error {

	// vérifier si il est vide
	if len(prefs) == 0 || len(alts) == 0 {
		return errors.New("Le profile list est NULL!")
	}

	Number_Alternative := len(alts)

	// Déterminer si tous les candidats n'apparaissent qu'une seule fois
	note := make(map[Alternative]bool)
	for it := range alts {
		if note[alts[it]] == false {
			note[alts[it]] = true
		} else{
			return errors.New("Le Alternative list est NULL!")
		}
	}

	// Juger tour à tour le résultat de chaque vote
	for i := 0; i < len(prefs); i++ {
		// Si la longueur est erronée, retournez directement
		if len(prefs[i]) != Number_Alternative {
			return errors.New("Profils ne sont pas tous complets!")
		}

		// Enregistrer tous les candidats dans un Set
		set := make(map[Alternative]int)
		for alt := range alts {
			set[alts[alt]] = 1
		}

		// Supprimer les candidats qui apparaissent immédiatement
		for j := 0; j < len(prefs[i]); j++ {
			set[prefs[i][j]]--
		}

		// Les éléments dans le set finalement doivent tous être 0, sinon, donne un erreur
		for _,j := range set {
			if j != 0 {
				return errors.New("Profils ne sont pas tous complets!");
			}
		}
	}

	return nil
}