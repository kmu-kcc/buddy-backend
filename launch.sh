export MONGO_URI="mongodb+srv://kmu-kcc:kmukcc1234@club.dvmoc.mongodb.net/myFirstDatabase?retryWrites=true&w=majority"
export ACCESS_SECRET="eprgjosefl"
make build
nohup ./buddy --port 3000 &