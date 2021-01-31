all: sass update_deps

sass: 
	sass app/sass/bulma.scss:public/css/bulma.css

update_deps:
	cp node_modules/jquery/dist/jquery.min.js public/js/
	cp node_modules/@fortawesome/fontawesome-free/js/all.min.js public/js/font-awesome.all.min.js
	cp node_modules/dropzone/dist/min/dropzone.min.js public/js/
	cp node_modules/dropzone/dist/min/dropzone.min.css public/css/
	cp node_modules/glightbox/dist/css/*.css public/css/
	cp node_modules/glightbox/dist/js/*.js public/js/
	cp -r node_modules/nanogallery2/dist/css/font public/css/
	cp node_modules/nanogallery2/dist/css/*.css public/css/
	cp node_modules/nanogallery2/dist/*.js public/js/
