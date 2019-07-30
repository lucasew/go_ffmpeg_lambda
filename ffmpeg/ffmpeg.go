package ffmpeg

import (
    "os/exec"
    "os"
    "fmt"
    "io"
)

type FFMpegSession struct {
    From io.Reader;
    To string;
    Params []string;
} // Caminho do arquivo, a parte de download rola separado

func (s FFMpegSession) Run() error {
    var params []string
    params = append(params, s.To)
    params = append(params, "-i")
    params = append(params, "pipe:")
    for _, v := range(s.Params) {
        params = append(params, v)
    }
    cmd := exec.Command("ffmpeg", params...)
    fmt.Printf("%s %v\n", "ffmpeg", params)
    cmd.Stdin = s.From
    cmd.Stdout = os.Stdout
    cmd.Stderr = os.Stderr
    return cmd.Run()
}
/*
func main() {
    s := FFMpegSession{
        from: "teste.png",
        to: "teste.jpg",
        params: []string{"-vf", "scale=800:-1"},
    }
    err := s.Run()
    if err != nil {
        panic(err)
    }
}
*/
