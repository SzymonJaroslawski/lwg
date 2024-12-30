package games

import (
	"io/fs"
	"os"
	"path/filepath"

	"github.com/SzymonJaroslawski/lwg/internal/utils"
	"github.com/google/uuid"
	"gopkg.in/yaml.v3"
)

type Games struct {
	games map[uuid.UUID]*Game
}

func (g *Games) AddGame(ga *Game, gameDir string) error {
	if g.games == nil {
		g.games = make(map[uuid.UUID]*Game)
	}

	if err := ga.SaveGameToDisk(gameDir); err != nil {
		return err
	}

	g.games[ga.ID] = ga
	return nil
}

func LoadGames(gameDir string) *Games {
	ga := &Games{}
	ga.games = make(map[uuid.UUID]*Game)

	filepath.WalkDir(gameDir, func(path string, d fs.DirEntry, err error) error {
		if d.IsDir() {
			return nil
		}

		f, err := os.Open(path)
		if err != nil {
			return err
		}
		defer f.Close()

		var g Game

		decoder := yaml.NewDecoder(f)
		err = decoder.Decode(&g)
		if err != nil {
			return err
		}

		ga.games[g.ID] = &g
		return nil
	})

	return ga
}

type Game struct {
	Name            string    `yaml:"name"`
	ExcecutablePath string    `yaml:"excecutable_path"`
	WinePrefixPath  string    `yaml:"wine_prefix_path"`
	ExtraArgs       string    `yaml:"extra_args"`
	ID              uuid.UUID `yaml:"id"`
	RunnerId        uuid.UUID `yaml:"runner_id"`
}

// Accepts Game struct, but will replace game.ID with random UUID.
// !! Manualy add it to the Games struct with Games.AddGame() !!
func NewGame(g Game) *Game {
	return &Game{
		Name:            g.Name,
		ExcecutablePath: g.ExcecutablePath,
		WinePrefixPath:  g.WinePrefixPath,
		ExtraArgs:       g.ExtraArgs,
		ID:              uuid.New(),
		RunnerId:        g.RunnerId,
	}
}

func (g *Game) SaveGameToDisk(path string) error {
	fileName := filepath.Join(path, g.Name+"_"+g.ID.String())

	if utils.Exsits(fileName) {
		err := os.Remove(fileName)
		if err != nil {
			return err
		}
	}

	f, err := os.Create(fileName)
	if err != nil {
		return err
	}
	defer f.Close()

	encoder := yaml.NewEncoder(f)
	err = encoder.Encode(g)
	if err != nil {
		return err
	}

	return nil
}
