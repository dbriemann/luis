sass: 
	sass app/sass/bulma.scss:public/css/bulma.css

update_deps:
	cp node_modules/jquery/dist/jquery.min.js public/js/
	cp node_modules/@fortawesome/fontawesome-free/js/all.min.js public/js/font-awesome.all.min.js
	cp node_modules/dropzone/dist/min/dropzone.min.js public/js/
	cp node_modules/dropzone/dist/min/dropzone.min.css public/css/
	cp node_modules/photoswipe/dist/photoswipe.min.js public/js/
	cp node_modules/photoswipe/dist/photoswipe.js public/js/
	cp node_modules/photoswipe/dist/photoswipe-ui-default.min.js public/js/
	cp node_modules/photoswipe/dist/photoswipe.css public/css/
	cp -r node_modules/photoswipe/dist/default-skin public/css/

prebuild: sass update_deps
