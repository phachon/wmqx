<?php
/**
 * wmqx publish
 * @author: phachon@163.com
 */

$wmqxPublishUrl = "http://127.0.0.1:3303/publish/";

// message info
$messageName = "ada";
$messageTokenHeader = "WMQX_MESSAGE_TOKEN";
$messageToken = "this is tokenssss";
$messageRouteKeyHeader = "WMQX_MESSAGE_ROUTEKEY";
$routeKey = "test222";

$url = $wmqxPublishUrl.$messageName;

// set http header
$headers = [
	$messageTokenHeader.":".$messageToken,
	$messageRouteKeyHeader.":".$routeKey,
];

// send value
$data = [
	"name" => "wmqx",
	"func" => "publish",
];
// if you want to send json data
$data = json_encode($data, true);

// start curl request
$ch = curl_init();
curl_setopt($ch, CURLOPT_URL, $url);
curl_setopt($ch, CURLOPT_HTTPHEADER, $headers);
curl_setopt($ch, CURLOPT_POST, TRUE);
curl_setopt($ch, CURLOPT_POSTFIELDS, $data);
curl_setopt($ch, CURLOPT_RETURNTRANSFER, TRUE);
$response = curl_exec($ch);
$httpCode = curl_getinfo($ch, CURLINFO_HTTP_CODE);
curl_close($ch);

var_dump($response);

//"{"code":1,"message":"success","data":{}}"