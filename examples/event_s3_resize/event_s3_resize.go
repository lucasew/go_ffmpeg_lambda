package main

import (
    "github.com/aws/aws-sdk-go/service/s3" 
    "github.com/aws/aws-sdk-go/service/s3/s3manager" 
    "github.com/aws/aws-sdk-go/aws/session" 
    "github.com/aws/aws-lambda-go/lambda"
    "github.com/aws/aws-lambda-go/events"
    "github.com/lucasew/go_ffmpeg_lambda/ffmpeg"
    "github.com/lucasew/golog"
    "os"
    "context"
    "fmt"
    "strconv"
)
// ======================= CONFIGURACOES =================
var ( // Nomes das variáveis de ambiente onde pegarei as configurações
    ENV_RESIZE_X = "CFG_RESIZE_X";
    ENV_RESIZE_Y = "CFG_RESIZE_Y";
    ENV_DESTINATION_BUCKET = "CFG_DESTINATION_BUCKET";
)

var log = golog.Default

type Request struct {
    Records []events.S3EventRecord `json:"Records"`;
}

var RESIZE_X int
var RESIZE_Y int
var DESTINATION_BUCKET string

func HandleRequest(ctx context.Context, ev Request) (string, error){
    logger := log.NewLogger(ev.Records[0].S3.Object.Key)
    sess := session.Must(session.NewSession())
    downloader := s3manager.NewDownloader(sess)
    f, err := os.Open("dummy.dat")
    if (err != nil) {
        return err.Error(), err
    }
    size, err := downloader.Download(f, &s3.GetObjectInput{
        Bucket: &ev.Records[0].S3.Bucket.Name,
        Key: &ev.Records[0].S3.Object.Key,
    })
    if (err != nil) {
        logger.Error(err.Error())
        return fmt.Sprintf("Cant write %d bytes to disk\n", size), err
    }
    ffsess := ffmpeg.FFMpegSession{
        From: f,
        To: ev.Records[0].S3.Object.Key,
        Params: []string{"-vf", fmt.Sprintf("scale=%d:%d", RESIZE_X, RESIZE_Y)},
    }
    err = ffsess.Run()
    if err != nil {
        logger.Error("Erro ffmpeg: %s", err.Error())
        return "erro ffmpeg", err
    }
    uploader := s3manager.NewUploader(sess)
    out, err := os.Open(ev.Records[0].S3.Object.Key)
    if err != nil {
        logger.Error("Erro ao abrir arquivo de saida: %s", err.Error())
        return "erro ao abrir arquivo de saida", err
    }
    _, err = uploader.Upload(&s3manager.UploadInput{
        Body: out,
        Bucket: &DESTINATION_BUCKET,
        Key: &ev.Records[0].S3.Object.Key,
    })
    if (err != nil) {
        logger.Error("Erro no upload: %s", err.Error())
        return "erro upload", err
    }
    return "", nil
}

func main() {
    log.Info("Iniciando...")
    var err error
    RESIZE_X, err = strconv.Atoi(os.Getenv(ENV_RESIZE_X))
    if err != nil {
        log.Warn("Não foi possível parsear número da variável %s, utilizando -1", ENV_RESIZE_X)
        RESIZE_X = -1
    }
    err = nil
    RESIZE_Y, err = strconv.Atoi(os.Getenv(ENV_RESIZE_Y))
    if err != nil {
        log.Warn("Não foi possível parsear número da variável %s, utilizando -1", ENV_RESIZE_Y)
        RESIZE_Y = -1
    }
    if (RESIZE_X == -1 && RESIZE_Y == -1) {
        panic(log.Error("Não é possível iniciar o sistema pois a altura e a largura do redimensionamento estão para automático (-1). Saindo..."))
    }
    DESTINATION_BUCKET = os.Getenv(ENV_DESTINATION_BUCKET)
    lambda.Start(HandleRequest)
}
