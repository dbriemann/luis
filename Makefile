sass: 
	sass app/sass/bulma.scss:public/css/bulma.css

update_deps:
	cp node_modules/jquery/dist/jquery.min.js public/js/jquery.min.js
	cp node_modules/@fortawesome/fontawesome-free/js/all.min.js public/js/font-awesome.all.min.js
	cp node_modules/dropzone/dist/min/dropzone.min.js public/js/dropzone.min.js
	cp node_modules/dropzone/dist/min/dropzone.min.css public/css/dropzone.min.css
	cp node_modules/dropzone/dist/basic.css public/css/basic.css

prebuild: sass update_deps
