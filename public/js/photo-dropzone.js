Dropzone.options.photoDropzone = {
  paramName: "file", // The name that will be used to transfer the file
  maxThumbnailFilesize: 20, // MB
  maxFilesize: 1000, // MB
  accept: function (file, done) {
    done();
  },
};
