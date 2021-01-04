sass: 
	sass app/sass/bulma.scss:public/css/bulma.css

update_deps:
	cp node_modules/jquery/dist/jquery.min.js public/js/jquery.min.js

prebuild: sass update_deps
