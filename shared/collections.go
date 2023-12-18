package shared

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/expr-lang/expr"
	"github.com/gobwas/glob"
	ortfodb "github.com/ortfo/db"
)

type HTMLString = string
type URLString = string
type AntonmedvExpression = string

func (c Collection) Works(db ortfodb.Database, tags []Tag, technologies []Technology) (worksInCollection []ortfodb.AnalyzedWork) {
	works := db.Works()
	for _, work := range works {
		contained, err := c.Contains(work, keys(works), tags, technologies)
		if err != nil {
			panic(fmt.Errorf("while checking if %s is in %v: %w", work.ID, c.Title, err))
		}
		if contained {
			worksInCollection = append(worksInCollection, work)
		}
	}
	return
}

func (c Collection) Contains(work ortfodb.AnalyzedWork, workIDs []string, tags []Tag, technologies []Technology) (bool, error) {
	context := map[string]interface{}{"work": work}
	for _, id := range workIDs {
		context[strings.ReplaceAll(id, "-", "_")] = id == work.ID
	}
	for _, t := range tags {
		isOnW := false
		for _, name := range work.Metadata.Tags {
			if t.ReferredToBy(name) {
				isOnW = true
				break
			}
		}
		context["tag_"+strings.ReplaceAll(t.URLName(), "-", "_")] = isOnW
	}
	for _, t := range technologies {
		isOnW := false
		for _, name := range work.Metadata.MadeWith {
			if t.ReferredToBy(name) {
				isOnW = true
				break
			}
		}
		context["technology_"+strings.ReplaceAll(t.Slug, "-", "_")] = isOnW
	}
	predicate, err := preprocessContainsPredicate(c.Includes, keys(context))
	if err != nil {
		return false, fmt.Errorf("while pre-processing work collection predicate: %w", err)
	}

	result, err := evaluateContainsPredicate(predicate, context)
	if err != nil {
		return false, fmt.Errorf("while evaluating work collection predicate %q (from %q): %w", predicate, c.Includes, err)
	}

	return result, nil
}

func preprocessContainsPredicate(expr AntonmedvExpression, variables []string) (AntonmedvExpression, error) {
	expr = regexp.MustCompile(`(\S)-(\S)`).ReplaceAllString(expr, "${1}_${2}")
	expr = regexp.MustCompile(`(\s|^)#(\S+)(\s|$)`).ReplaceAllString(expr, "${1}tag_${2}${3}")
	expr = regexp.MustCompile(`(\s|^)made with (\S+)(\s|$)`).ReplaceAllString(expr, "${1}technology_${2}${3}")
	for _, captures := range regexp.MustCompile(`(\s|^)((?:.*\*.*)+)(\s|$)`).FindAllStringSubmatch(expr, -1) {
		globPattern, err := glob.Compile(captures[2])
		if err != nil {
			return "", fmt.Errorf("while compiling glob pattern %s: %w", captures[2], err)
		}

		matchingVariables := make([]string, 0)
		for _, variable := range variables {
			if globPattern.Match(variable) {
				matchingVariables = append(matchingVariables, variable)
			}
		}

		expr = strings.ReplaceAll(expr, captures[0], captures[1]+"("+strings.Join(matchingVariables, " or ")+")"+captures[3])
	}
	return expr, nil
}

func evaluateContainsPredicate(preprocessedExpr AntonmedvExpression, context map[string]interface{}) (bool, error) {
	compiledExpr, err := expr.Compile(preprocessedExpr)
	if err != nil {
		return false, fmt.Errorf("invalid work collection predicate: %w", err)
	}

	value, err := expr.Run(compiledExpr, context)
	if err != nil {
		return false, fmt.Errorf("couldn't evaluate predicate: %w", err)
	}

	switch coerced := value.(type) {
	case bool:
		return coerced, nil
	default:
		return false, fmt.Errorf("predicate does not evaluate to a boolean, but to %#v", value)
	}
}

func keys[K comparable, V any](m map[K]V) []K {
	keys := make([]K, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	return keys
}
