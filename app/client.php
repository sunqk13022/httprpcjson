<?php
$uri = "http://127.0.0.1:1234/jsonrpc/";
$adb_handle = curl_init();
$json_array["jsonrpc"] = "2.0";

$json_array["method"] = "Counter.Get";
$args = array();
$args["A"] = 2;
$args["B"] = 12;
$json_array["params"] = $args;
$json_array["id"] = '2';

$postdata= json_encode($json_array);

$adb_handle = curl_init($uri);
curl_setopt($adb_handle, CURLOPT_POSTFIELDS, $postdata);
curl_setopt($adb_handle, CURLOPT_RETURNTRANSFER, true);
curl_setopt($adb_handle, CURLOPT_HTTPHEADER, array(
  'Content-Type: application/json',
  'Content-Length: ' . strlen($postdata))
);
$responce =  curl_exec($adb_handle);
curl_close($adb_handle);
echo $responce . "\n";
?>
