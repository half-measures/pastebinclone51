package ui

//can only embed files within same directory
import "embed"

//below is not a comment, its comment directive, instructs go to store files from ui/html and static folders into a embed.FS filesystem

//go:embed "html" static
var Files embed.FS
