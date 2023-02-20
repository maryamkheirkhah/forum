echo "Building image"
docker build -t  forum-app . 
echo "Creating a container from the ascii-art-app image"
docker run -p 8080:8080 -d --name gritforum forum-app
echo "Running container list"
docker ps
echo "Image list"
docker image list
echo "Container list"
docker ps -a