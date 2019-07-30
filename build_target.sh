wget https://johnvansickle.com/ffmpeg/releases/ffmpeg-release-amd64-static.tar.xz -O out.tar.xz
tar -xvf out.tar.xz
rm out.tar.xz
mv ffmpeg-* downloaded
mkdir -p target
mv downloaded/ffmpeg target
rm downloaded -rf
go build -o target/main ./examples/event_s3_resize/event_s3_resize.go 

zip deploy_me.zip target/*
