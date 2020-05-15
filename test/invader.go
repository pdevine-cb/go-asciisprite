package main

import (
	"fmt"
	"math"
	"math/rand"
	"time"

	sprite "github.com/pdevine/go-asciisprite"
	tm "github.com/pdevine/go-asciisprite/termbox"
)

var allSprites sprite.SpriteGroup
var Width int
var Height int
var gameState *GameState

var randSrc *rand.Rand


type EdgeType int

const (
	UpperLeftEdge EdgeType = iota
	UpperRightEdge
	LowerLeftEdge
	LowerRightEdge
)

type GameStateScreen int

const (
	Title GameStateScreen = iota
	Play
	GameOver
)

type GameState struct {
	invaders  []*Invader
	direction int
	nextFire  *NextInvaderToFire
	player    *Fighter
	lives     []*Fighter
	State     GameStateScreen
	Score     *Score
	UfoTimer  int
	Ufo       *Ufo
}

type Invader struct {
	sprite.BaseSprite
	Timer     int
	TimeOut   int
	Type	  int
	Col	  int
	Exploding bool
	Dead      bool
}

type NextInvaderToFire struct {
	Invader *Invader
	Timer   int
	TimeOut int
	Fired	bool
}

type Ufo struct {
	sprite.BaseSprite
	Timer     int
	TimeOut   int
	Direction int
	Exploding bool
	Dead      bool
}

type Fighter struct {
	sprite.BaseSprite
	Timer     int
	TimeOut   int
	VX        float64
	AX        float64
	Counter   int
	Exploding bool
	Dead      bool
}

type Bullet struct {
	sprite.BaseSprite
	VY	int
	Dead	bool
}

type Score struct {
	sprite.BaseSprite
	Val int
}

type Logo struct {
	sprite.BaseSprite
	TargetY int
	VY      float64
	DY      float64
}

type Arrow struct {
	sprite.BaseSprite
	DX    float64
	DY    float64
	Angle float64
	Type  EdgeType
}

type Edge struct {
	sprite.BaseSprite
}

const logo = `                                                                                      
                               XXXX                                                    
                             XXXXX   XXXXXXX          XXXXXXXXX                        
                           XXXXX    XXXXXXX          XXXXXXXXX                         
                         XXX XXX   XXX     XXX   XXX   XXX                             
                             XXX  XXXXX     XXX XXX   XXX                              
                             XXX  XX         XXX     XXX                               
                             XXX  XXXXX   XXX  XXX  XXX                                
                                                                                       
   XXXXXXXXX  XXX  XXX  XXX     XXX    XXXX    XXXXXXX      XXXXXX  XXXXXX     XXXXXX  
    XXXXXXXXX  XXX  XXX  XXX    XXX   XX  XX   XXX   XX    XXXXXX  XX  XXX    XXXXXX   
        XXX     XXXXXXXX  XXX   XXX   XX  XX   XXX   XX    XX     XX   XXX  XXX        
         XXX     XXXXXXXX  XXX  XXX  XX    XX  XXX  XX   XXXX     XX  XXX    XXX       
          XXX     XXX XXXX  XXX XXX  XXXXXXXX  XXX  XX   XXXX    XXXXXX       XXX      
           XXX     XXX  XXX  XXXXXX  XXXXXXXX  XXX  XX  XX       XX  XX         XXX    
         XXXXXXXXX  XXX  XXX  XXXXX  XX    XX  XXXXXX  XXXXX   XXX  XX     XXXXXX      
          XXXXXXXXX  XXX  XXX  XXX  XXX    XXX XXXXXX  XXXXXX  XXX   XXX  XXXXXX       
                                                                                       `

/*
const invader_c0 = `  X     X
 XXXXXXXXX
 XX XXX XX
 XXXXXXXXX
 X X   X X
X X   X X`

const invader_c1 = `  X     X
 XXXXXXXXX
 XX XXX XX
 XXXXXXXXX
 X X   X X 
  X X   X X`
*/

const invader_c0 = `  X     X   
   X   X    
  XXXXXXX   
 XX XXX XX  
XXXXXXXXXXX 
X XXXXXXX X 
X X     X X 
   XX XX     `

const invader_c1 = `  X     X   
X  X   X  X 
X XXXXXXX X 
XXX XXX XXX 
XXXXXXXXXXX 
 XXXXXXXXX  
  X     X   
 X       X  `


const invader2_c0 = `   XX     
  XXXX   
 XXXXXX  
XX XX XX 
XXXXXXXX 
  X  X   
 X XX X  
X X  X X  `

const invader2_c1 = `   XX    
  XXXX   
 XXXXXX  
XX XX XX 
XXXXXXXX 
  X  X   
 X    X  
  X  X    `

const invader3_c0 = `    XXXX   
 XXXXXXXXXX  
XXXXXXXXXXXX 
XXX  XX  XXX 
XXXXXXXXXXXX 
   XX  XX    
  X  XX  X   
   X    X    `

const invader3_c1 = `    XXXX   
 XXXXXXXXXX  
XXXXXXXXXXXX 
XXX  XX  XXX 
XXXXXXXXXXXX 
   XX  XX    
  XX XX XX   
XX        XX `


const fighter_c0 = `
      X       
     XXX      
     XXX     
 XXXXXXXXXXX  
XXXXXXXXXXXXX
XXXXXXXXXXXXX
XXXXXXXXXXXXX`


const fighter_explode_c0 = `
   X
     X X X   
   X X       
      XX XX  
X   X XX X   
  XXXXXXXX X 
 XXXXXXXXXX X`

const fighter_explode_c1 = `

             
    X        
  X    X X X 
X   X XX     
  X XX X X X 
 XX XXX XXX X`


const explosion_c1 = `X  X   X  X  
 X  X X  X   
  X     X    
XX       XX  
  X     X    
 X  X X  X   
X  X   X  X 
            `

const explosion_c0 = `        
         
     XX  
   X    X
   X    X
     XX  `


const ufo_c0 = `
     XXXXXXX      
   XXXXXXXXXXX    
  XXXXXXXXXXXXX   
 X XXX XXX XXX X  
XXXXXXXXXXXXXXXXX 
  XXX  XXX  XXX   
   X         X     `

const ufo_c1 = `
     XXXXXXX      
   XXXXXXXXXXX    
  XXXXXXXXXXXXX   
 XXX XXX XXX XXX  
XXXXXXXXXXXXXXXXX 
  XXX  XXX  XXX   
   X         X     `


const frame_ul = `XXXXXXXX
XXXXXXXX
XX      
XX      
XX      
XX      `

const frame_ur = `XXXXXXXX
XXXXXXXX
      XX
      XX
      XX
      XX`

const frame_ll = `XX      
XX      
XX      
XX      
XXXXXXXX
XXXXXXXX`

const frame_lr = `      XX
      XX
      XX
      XX
XXXXXXXX
XXXXXXXX`

const arrow_ul = `      
XXXXX  
XXX    
X XX   
X  XX  
    XX `

const arrow_ur = `        
   XXXXX
     XXX
    XX X
   XX  X
  XX    `

const arrow_ll = `        
    XX  
X  XX   
X XX    
XXX     
XXXXX   `

const arrow_lr = `        
  XX    
   XX  X
    XX X
     XXX
   XXXXX`

func NewGame() *GameState {
	gs := &GameState{State: Title}
	return gs
}

func (gs *GameState) StartGame() {
	allSprites.RemoveAll()

	gs.player = NewFighter()
	gs.player.X = 34
	allSprites.Sprites = append(allSprites.Sprites, gs.player)

	gs.Score = &Score{BaseSprite: sprite.BaseSprite{
		Visible: true,
		X:       20,
		Y:       1},
	}

	gs.Score.AddCostume(sprite.Convert(sprite.BuildString(fmt.Sprintf("%06d", gs.Score.Val))))
	allSprites.Sprites = append(allSprites.Sprites, gs.Score)

	for i := 0; i < 2; i++ {
		f := NewFighter()
		f.Y = 75
		f.X = i*8 + 2
		gs.lives = append(gs.lives, f)
		allSprites.Sprites = append(allSprites.Sprites, f)
	}
	gs.State = Play
}

func (gs *GameState) Update() {
	t := 150
	gs.cullSprites()
	if gs.State == Play || gs.State == GameOver {
		gs.checkDirection()
		gs.nextFire.Update()
		if gs.nextFire.Fired {
			gs.setNextInvaderToFire()
		}
		t = 500
	}
	gs.UfoTimer++
	if gs.UfoTimer > t {
		gs.UfoTimer = 0
		gs.Ufo = NewUfo()
		if gs.State == Title {
			gs.Ufo.Y = 34
		}
		allSprites.Sprites = append(allSprites.Sprites, gs.Ufo)
	}
}

func (gs *GameState) setNextInvaderToFire() {
	var invaders map[int]int
	invaders = make(map[int]int)
	cols := []int{}

	// iterate through each of the invaders and determine which columns they are in
	for _, i := range gs.invaders {
		if !i.Dead && !i.Exploding {
			invaders[i.Col] += 1
		}
	}

	// get the keys of the columns and choose one at random
	for k, _ := range invaders {
		cols = append(cols, k)
	}
	if len(cols) == 0 {
		return
	}
	i_col := randSrc.Intn(len(cols))

	// iterate backwards through the invaders and have the first one we find fire
	for i := len(gs.invaders) - 1; i >= 0; i-- {
		if gs.invaders[i].Col == cols[i_col] && !gs.invaders[i].Dead && !gs.invaders[i].Exploding {
			gs.nextFire.Invader = gs.invaders[i]
			break
		}
	}

	gs.nextFire.Timer = 0
	gs.nextFire.TimeOut = 5
	gs.nextFire.Fired = false
}

func (gs *GameState) checkDirection() {
	changeDir := false
	dead := false
	var d int
	for _, s := range gs.invaders {
		if s.X > 98 && gs.direction > 0 {
			d = -(s.X - 98)
			changeDir = true
			break
		} else if s.X < 3  && gs.direction < 0 {
			d = 3 - s.X
			changeDir = true
			break
		}
	}
	if changeDir {
		for _, s := range gs.invaders {
			s.X += d
			if gs.State != GameOver {
				s.Y += 3
				if gs.player.Y - s.Y < 5 {
					dead = true
				}
			}
		}
		gs.direction = -gs.direction
	}
	if dead {
		gs.player.Explode()
		for _, f := range gs.lives {
			f.Explode()
		}
		gs.State = GameOver
	}
}

func (gs *GameState) cullSprites() {
	for _, s := range allSprites.Sprites {
		switch sp := s.(type) {
		case *Bullet:
			if sp.Dead {
				allSprites.Remove(sp)
			}
		case *Invader:
			if sp.Dead {
				gs.removeInvader(sp)
				allSprites.Remove(sp)
				l := len(gameState.invaders)-1
				speedUps := []int{36, 27, 18, 9, 3, 2}
				for _, sp := range speedUps {
					if l == sp {
						for _, i := range gameState.invaders {
							i.TimeOut--
						}
					}
				}
			}
		case *Ufo:
			if sp.Dead {
				allSprites.Remove(sp)
				gs.Ufo = nil
			}
		case *Fighter:
			if sp.Dead {
				allSprites.Remove(sp)
				if len(gs.lives) == 0 {
					gs.State = GameOver
				} else {
					gs.player = gs.lives[len(gs.lives)-1]
					gs.lives = gs.lives[:len(gs.lives)-1]
					gs.player.Y = 70
				}
			}
		default:
		}
	}
}

func (gs *GameState) removeInvader(i *Invader) {
	var idx int
	for cnt, s := range gs.invaders {
		if i == s {
			idx = cnt
			break
		}
	}
	copy(gs.invaders[idx:], gs.invaders[idx+1:])
	gs.invaders[len(gs.invaders)-1] = nil
	gs.invaders = gs.invaders[:len(gs.invaders)-1]
}

func (gs *GameState) createWave() {
	gs.direction = 2

	for cnt := 0; cnt < 9; cnt++ {
		i := NewInvader(0)
		i.Col = cnt
		i.X = cnt * 10 + 10
		i.Y = 7
		gs.invaders = append(gs.invaders, i)
		allSprites.Sprites = append(allSprites.Sprites, i)
	}

	for y := 0; y < 2; y++ {
		for cnt := 0; cnt < 9; cnt++ {
			i := NewInvader(1)
			i.Col = cnt
			i.X = cnt * 10 + 9
			i.Y = y * 6 + 13
			gs.invaders = append(gs.invaders, i)
			allSprites.Sprites = append(allSprites.Sprites, i)
		}
	}

	for y := 0; y < 2; y++ {
		for cnt := 0; cnt < 9; cnt++ {
			i := NewInvader(2)
			i.Col = cnt
			i.X = cnt * 10 + 9
			i.Y = y * 6 + 25
			gs.invaders = append(gs.invaders, i)
			allSprites.Sprites = append(allSprites.Sprites, i)
		}
	}

	gs.nextFire = &NextInvaderToFire{
		Invader: gs.invaders[len(gs.invaders)-1],
		TimeOut: 5,
		Fired:   false,
	}
}

func (n *NextInvaderToFire) Update() {
	n.Timer = n.Timer + 1
	if n.Timer > n.TimeOut {
		n.Invader.Fire()
		n.Timer = 0
		n.Fired = true
	}
}

func NewInvader(t int) *Invader {
	s := &Invader{BaseSprite: sprite.BaseSprite{
		Visible:        true,
		X:		2,
		Y:              2,
		CurrentCostume: 0},
		Type:    t,
		Timer:   0,
		TimeOut: 6}

	if t == 1 {
		s.AddCostume(sprite.Convert(invader_c0))
		s.AddCostume(sprite.Convert(invader_c1))
	} else if t == 0 {
		s.AddCostume(sprite.Convert(invader2_c0))
		s.AddCostume(sprite.Convert(invader2_c1))
	} else if t == 2 {
		s.AddCostume(sprite.Convert(invader3_c0))
		s.AddCostume(sprite.Convert(invader3_c1))
	}
	return s
}

func (s *Invader) Explode() {
	s.Costumes = []*sprite.Costume{}
	s.AddCostume(sprite.Convert(explosion_c0))
	s.AddCostume(sprite.Convert(explosion_c1))
	s.Exploding = true
}

func (s *Invader) Update() {
	s.Timer = s.Timer + 1
	if s.Timer > s.TimeOut {
		s.X = s.X + gameState.direction
		s.CurrentCostume = s.CurrentCostume + 1
		if s.CurrentCostume >= len(s.Costumes) {
			s.CurrentCostume = 0
			if s.Exploding {
				s.Dead = true
			}
		}
		s.Timer = 0
	}
}

func (s *Invader) Fire() {
	b := &Bullet{BaseSprite: sprite.BaseSprite{
		Visible: true,
		X: s.X+3,
		Y: s.Y+4},
		VY: 2}
	b.AddCostume(sprite.Convert("X \n X\nX \n  "))
	b.AddCostume(sprite.Convert(" X\nX \n X\n  "))
	allSprites.Sprites = append(allSprites.Sprites, b)
}

func NewUfo() *Ufo {
	var x int
	ds := []int{-1, 1}
	d := ds[randSrc.Intn(len(ds))]
	if d < 0 {
		x = 100
	} else {
		x = -4
	}
	s := &Ufo{BaseSprite: sprite.BaseSprite{
		Visible: true,
		X:       x,
		Y:       4},
		Direction: d,
		TimeOut:   3}
	s.AddCostume(sprite.Convert(ufo_c0))
	s.AddCostume(sprite.Convert(ufo_c1))
	return s
}

func (s *Ufo) Update() {
	if s.Direction > 0 {
		s.X += 1
		if s.X > 100 {
			s.Dead = true
		}
	} else {
		s.X -= 1
		if s.X+s.Width < 0 {
			s.Dead = true
		}
	}
	s.Timer = s.Timer + 1
	if s.Timer > s.TimeOut {
		s.CurrentCostume = s.CurrentCostume + 1
		if s.CurrentCostume >= len(s.Costumes) {
			s.CurrentCostume = 0
			if s.Exploding {
				s.Dead = true
			}
		}
		s.Timer = 0
	}
}

func (s *Ufo) Explode() {
	s.Costumes = []*sprite.Costume{}
	s.AddCostume(sprite.Convert(explosion_c0))
	s.AddCostume(sprite.Convert(explosion_c1))
	s.Exploding = true
	scores := []int{50, 100, 150, 200, 250, 300}
	score := randSrc.Intn(len(scores))
	gameState.Score.Val += scores[score]

}

func NewFighter() *Fighter {
	s := &Fighter{BaseSprite: sprite.BaseSprite{
		Visible:	true,
		X: 2,
		Y: 70},
		TimeOut: 2,
		}
	s.AddCostume(sprite.Convert(fighter_c0))
	return s
}

func (s *Fighter) MoveLeft() {
	if !s.Exploding || !s.Dead {
		s.AX = -3
		s.VX = 0
	}
}

func (s *Fighter) MoveRight() {
	if !s.Exploding || !s.Dead {
		s.AX = 3
		s.VX = 0
	}
}

func (s *Fighter) Fire() {
	if s.Exploding || s.Dead {
		return
	}

	b := &Bullet{BaseSprite: sprite.BaseSprite{
		Visible: true,
		X: s.X+3,
		Y: s.Y},
		VY: -2}
	b.AddCostume(sprite.Convert("X \n  "))
	allSprites.Sprites = append(allSprites.Sprites, b)
}

func (s *Fighter) Update() {
	s.VX = s.VX + s.AX
	s.AX = 0
	s.VX *= 0.85
	s.X += int(math.Round(s.VX))

	if s.X < 2 {
		s.X = 2
		s.VX = 0
		s.AX = 0
	} else if s.X > 100 {
		s.X = 100
		s.VX = 0
		s.AX = 0
	}

	if s.Exploding {
		s.Timer += 1
		if s.Timer > s.TimeOut {
			s.CurrentCostume = s.CurrentCostume + 1
			if s.CurrentCostume >= len(s.Costumes) {
				s.CurrentCostume = 0
			}
			s.Timer = 0
			s.Counter += 1
			if s.Counter > 4 {
				s.Dead = true
			}
		}
	}
}

func (s *Fighter) Explode() {
	s.Costumes = []*sprite.Costume{}
	s.AddCostume(sprite.Convert(fighter_explode_c0))
	s.AddCostume(sprite.Convert(fighter_explode_c1))
	s.Exploding = true
}

func (s *Bullet) Update() {
	s.Y += s.VY

	if s.Y < 0 || s.Y > 80 {
		s.Dead = true
		return
	}

	s.CurrentCostume = s.CurrentCostume + 1
	if s.CurrentCostume >= len(s.Costumes) {
		s.CurrentCostume = 0
	}

	if gameState.player.HitAtPoint(s.X, s.Y+1) {
		gameState.player.Explode()
	}

	for _, i := range gameState.invaders {
		if !i.Exploding && i.HitAtPoint(s.X, s.Y) {
			if i.Type == 0 {
				gameState.Score.Val += 30
			} else if i.Type == 1 {
				gameState.Score.Val += 20
			} else {
				gameState.Score.Val += 10
			}
			i.Explode()
			s.Dead = true
		}
	}

	if gameState.Ufo != nil && !gameState.Ufo.Exploding && gameState.Ufo.HitAtPoint(s.X, s.Y+1) {
		gameState.Ufo.Explode()
	}

}

func (s *Score) Update() {
	s.Costumes = nil
	s.AddCostume(sprite.Convert(sprite.BuildString(fmt.Sprintf("score %06d", s.Val))))
}

func (s *Logo) Update() {
	s.VY = (float64(s.TargetY) - float64(s.Y)) * 0.3
	s.Y += int(math.Round(s.VY))
}

func NewArrow(t EdgeType) *Arrow {
	s := &Arrow{BaseSprite: sprite.BaseSprite{
		Visible: true},
		Type: t,
	}
	switch s.Type {
	case UpperLeftEdge:
		s.DX = 4
		s.DY = 3
		s.AddCostume(sprite.Convert(arrow_ul))
	case UpperRightEdge:
		s.DX = 96
		s.DY = 3
		s.AddCostume(sprite.Convert(arrow_ur))
	case LowerLeftEdge:
		s.DX =  4
		s.DY = 74
		s.AddCostume(sprite.Convert(arrow_ll))
	case LowerRightEdge:
		s.DX = 96
		s.DY = 74
		s.AddCostume(sprite.Convert(arrow_lr))
	}
	return s
}

func (s *Arrow) Update() {
	s.Angle += 0.25

	d := math.Sin(s.Angle) * 0.2
	switch s.Type {
	case UpperLeftEdge:
		s.DX += d
		s.DY += d
	case UpperRightEdge:
		s.DX -= d
		s.DY += d
	case LowerLeftEdge:
		s.DX += d
		s.DY -= d
	case LowerRightEdge:
		s.DX -= d
		s.DY -= d
	}
	s.X = int(math.Round(s.DX))
	s.Y = int(math.Round(s.DY))

}

func NewEdge(t EdgeType) *Edge {
	s := &Edge{BaseSprite: sprite.BaseSprite{
		Visible: true},
	}

	switch t {
	case UpperLeftEdge:
		s.X = 1
		s.Y = 1
		s.AddCostume(sprite.Convert(frame_ul))
	case UpperRightEdge:
		s.X = 99
		s.Y = 1
		s.AddCostume(sprite.Convert(frame_ur))
	case LowerLeftEdge:
		s.X = 1
		s.Y = 76
		s.AddCostume(sprite.Convert(frame_ll))
	case LowerRightEdge:
		s.X = 99
		s.Y = 76
		s.AddCostume(sprite.Convert(frame_lr))
	}
	return s
}


func ShowTitle() {
	l := &Logo{BaseSprite: sprite.BaseSprite{
		X:       30,
		Y:       -10,
		Visible: true},
		TargetY: 10,
	}
	l.AddCostume(sprite.Convert(logo))

	adj_txt := &sprite.BaseSprite{
		X: 20,
		Y: 22,
		Visible: true,
	}
	adj_txt.AddCostume(sprite.NewCostume("ADJUST YOUR TERMINAL TO SEE ALL OF THE EDGES OF THE PLAY AREA", '@'))

	font_txt := &sprite.BaseSprite{
		X: 27,
		Y: 24,
		Visible: true,
	}
	font_txt.AddCostume(sprite.NewCostume("Recommended Font:  Menlo 8pt (0.81 Line Spacing)", '@'))

	for _, et := range []EdgeType{UpperLeftEdge, UpperRightEdge, LowerLeftEdge, LowerRightEdge} {
		a := NewArrow(et)
		e := NewEdge(et)
		allSprites.Sprites = append(allSprites.Sprites, e)
		allSprites.Sprites = append(allSprites.Sprites, a)

	}

	allSprites.Sprites = append(allSprites.Sprites, l)
	allSprites.Sprites = append(allSprites.Sprites, adj_txt)
	allSprites.Sprites = append(allSprites.Sprites, font_txt)
}

func main() {
	// XXX - Wait a bit until the terminal is properly initialized
	time.Sleep(500 * time.Millisecond)

	err := tm.Init()
	if err != nil {
		panic(err)
	}
	defer tm.Close()

	s := rand.NewSource(time.Now().Unix())
	randSrc = rand.New(s)

	Width, Height = tm.Size()

	event_queue := make(chan tm.Event)
	go func() {
		for {
			event_queue <- tm.PollEvent()
		}
	}()

	gameState = NewGame()

	ShowTitle()

mainloop:
	for {
		tm.Clear(tm.ColorDefault, tm.ColorDefault)

		select {
		case ev := <-event_queue:
			if ev.Type == tm.EventKey {
				if gameState.State == Title {
					gameState.StartGame()
					continue
				}

				if ev.Key == tm.KeyEsc || ev.Ch == 'q' {
					break mainloop
				} else if gameState.State == Play {
					if ev.Key == tm.KeyArrowLeft {
						gameState.player.MoveLeft()
					} else if ev.Key == tm.KeyArrowRight {
						gameState.player.MoveRight()
					} else if ev.Key == tm.KeySpace {
						gameState.player.Fire()
					}
				}
			} else if ev.Type == tm.EventResize {
				Width = ev.Width
				Height = ev.Height
			}
		default:
			if gameState.State == Play {
				if len(gameState.invaders) == 0 {
					gameState.createWave()
				}
			}
			allSprites.Update()
			gameState.Update()
			allSprites.Render()
			time.Sleep(50 * time.Millisecond)
		}
	}
}
