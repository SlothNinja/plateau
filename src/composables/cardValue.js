export function useCardValue(card) {
  let v = 0
  switch (card.Rank) {
    case "one":
      v += 1
    case "two":
      v += 2
    case "three":
      v += 3
    case "four":
      v += 4
    case "five":
      v += 5
    case "six":
      v += 6
    case "seven":
      v += 7
    case "eight":
      v += 8
    case "nine":
      v += 9
    case "ten":
      v += 10
    case "eleven":
      v += 11
    case "valet":
      v +=11
    case "twelve":
      v += 12
    case "cavalier":
      v += 12
    case "thirteen":
      v += 13
    case "dame":
      v += 13
    case "fourteen":
      v += 14
    case "roi":
      v += 14
    case "fifteen":
      v += 15
    case "sixteen":
      v += 16
    case "seventeen":
      v += 17
    case "eighteen":
      v += 18
    case "nineteen":
      v += 19
    case "twenty":
      v += 20
    case "twentyone":
      v += 21
    case "excuse":
      v += 22
    default:
      v += 0
  }

  if (card.suit == "trumps") {
    v += 100
  }

  return v
}
