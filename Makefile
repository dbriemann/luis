sass: 
	sass app/sass/bulma.scss:public/css/bulma.css

update_deps:
	cp node_modules/jquery/dist/jquery.min.js public/js/
	cp node_modules/@fortawesome/fontawesome-free/js/all.min.js public/js/font-awesome.all.min.js
	cp node_modules/dropzone/dist/min/dropzone.min.js public/js/
	cp node_modules/dropzone/dist/min/dropzone.min.css public/css/
	cp node_modules/lightgallery.js/dist/js/lightgallery.js public/js/
	cp node_modules/lightgallery.js/dist/js/lightgallery.min.js public/js/
	cp node_modules/lightgallery.js/dist/css/lightgallery.css public/css/
	cp node_modules/lightgallery.js/dist/css/lightgallery.min.css public/css/
	cp node_modules/lightgallery.js/dist/fonts/*.* public/fonts/
	cp node_modules/lightgallery.js/dist/img/*.* public/img/
	cp node_modules/lg-thumbnail/dist/lg-thumbnail.js public/js/
	cp node_modules/lg-thumbnail/dist/lg-thumbnail.min.js public/js/

prebuild: sass update_deps
