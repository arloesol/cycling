echo "moving gpx files"
cd gpx || exit 1
find . -type d -empty -delete
mv * $GITDIR/static/gpxfiles/
cd ..
echo "moving img files"
cd img || exit 2
find . -type d -empty -delete
mv gallery/* $GITDIR/assets/routes/gallery/
mkdir $GITDIR/assets/signage
mv signage/* $GITDIR/assets/signage/
cd ..
echo "moving route pages"
cd route
find . -type f -empty -delete
mv * $GITDIR/content/route/
cd ..
echo "done"
