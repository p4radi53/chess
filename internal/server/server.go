package server

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/p4radi53/chess/internal/chess"
)

var game *chess.Game

func setupCorsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Content-Type")
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	}
}

func StartServer() {
	game = chess.NewGame()
	r := gin.Default()
	r.Use(func(c *gin.Context) {
		setupCorsMiddleware()(c)
	})
	r.GET("/game", handleGetGame)
	r.POST("/move", handleMove)
	r.POST("/reset", handleNewGame)
	r.GET("/legal-moves", handleShowLegalMoves)
	r.GET("/ping", handlePing)

	if err := r.Run(); err != nil {
		log.Fatalf("server failed: %v", err)
	}
}

func handleGetGame(c *gin.Context) {
	c.JSON(200, game)
}

func handlePing(c *gin.Context) {
	c.JSON(200, gin.H{"message": "pong"})
}

func handleNewGame(c *gin.Context) {
	game = chess.NewGame()
	c.JSON(200, game)
}

func handleMove(c *gin.Context) {
	var move struct {
		FromFile int `json:"from_file"`
		FromRank int `json:"from_rank"`
		ToFile   int `json:"to_file"`
		ToRank   int `json:"to_rank"`
	}
	if err := c.BindJSON(&move); err != nil {
		c.JSON(400, gin.H{"error": "Invalid move format", "detail": err.Error()})
		return
	}
	err := game.MakeMove(move.FromFile, move.FromRank, move.ToFile, move.ToRank)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, game)
}

func handleShowLegalMoves(c *gin.Context) {
	var request struct {
		File int `form:"file"`
		Rank int `form:"rank"`
	}
	if err := c.BindQuery(&request); err != nil {
		c.JSON(400, gin.H{"error": "Invalid request format", "detail": err.Error()})
		return
	}
	moves := game.Board.LegalMoves(chess.Square{File: request.File, Rank: request.Rank}, game.LastMove())
	if moves == nil {
		moves = []chess.Square{}
	}
	c.JSON(200, moves)
}
