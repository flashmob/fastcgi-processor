<?php
# Email saving script example
# This does not need to be in a web root, it will be executed directly by php-fpm.
# Simply point the fcgi_script_filename_save config option to the path of this file
# (making sure that your fcgi has permissions)

$all = var_export($_POST, true);

file_put_contents(__DIR__ . "/saved", $all, FILE_APPEND);

// respond with SAVED if everything went OK
echo "SAVED";