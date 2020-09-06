
package main

func selectCard(params Params) int {
    result := 0

    if params.mode == MODE_MAX {
        result = maxCard(params.cards)
    } else if params.mode == MODE_MIN {
        result = minCard(params.cards, params.maxCard)
    }

    return result
}

func maxCard(cards []int) int {
	result := 0

	for _, card := range cards {
		if card > result {
			result = card
		}
	}

	return result
}

func minCard(cards []int, max int) int {
	result := max * 2

	for _, card := range cards {
		if card < result {
			result = card
		}
	}

	return result
}
