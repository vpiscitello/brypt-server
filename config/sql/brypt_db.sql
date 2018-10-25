DROP TABLE IF EXISTS `brypt_users`;

CREATE TABLE `brypt_users` (
	`user_id` char(32) NOT NULL,
	`username` varchar(32) NOT NULL,
	`first_name` varchar(32) NOT NULL,
	`last_name` varchar(255) NOT NULL,
	`email` varchar(255) NOT NULL,
	`organization` varchar(255),
	`age` datetime,
	`join_date` datetime,
	`last_login` datetime,
	`login_attempts` int(11),
	`login_token` char(255),
	`region` char(255),
	PRIMARY KEY (`user_id`) 
) ENGINE=InnoDB DEFAULT CHARSET=latin1;


DROP TABLE IF EXISTS `brypt_nodes`;

CREATE TABLE `brypt_nodes` (
	`node_id` char(32) NOT NULL,
	`serial_number` char(255) NOT NULL,
	`type` varchar(32),
	`created_on` datetime,
	`registered_on` datetime,
	`registered_to` char(32),
	`connected_network`char(32),
	PRIMARY KEY (`node_id`)
) ENGINE=InnoDB DEFAULT CHARSET=latin1;

DROP TABLE IF EXISTS `brypt_users_networks`;

CREATE TABLE `brypt_users_networks` (
	`user_id` char(32) NOT NULL,
	`network_id` char(32) NOT NULL,
	PRIMARY KEY (`user_id`, `network_id`),
	CONSTRAINT `users_networks_fk1` FOREIGN KEY (`user_id`) REFERENCES `brypt_users` (`user_id`)
) ENGINE=InnoDB DEFAULT CHARSET=latin1;
