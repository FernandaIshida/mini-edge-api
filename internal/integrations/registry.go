package integrations

var Registry = map[string]Integration{
	"github": {
		Name: "github",
		URL:  "https://api.github.com",
	},
	"pokemon": {
		Name: "pokemon",
		URL:  "https://pokeapi.co/api/v2/pokemon/pikachu",
	},
	"weather": {
		Name: "weather",
		URL:  "https://api.open-meteo.com/v1/forecast?latitude=0&longitude=0&current=temperature_2m",
	},
	"exchange": {
		Name: "exchange",
		URL:  "https://api.exchangerate.host/latest",
	},
}
