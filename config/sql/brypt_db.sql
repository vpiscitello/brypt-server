--Add create table methods here

DROP TABLE IF EXISTS `brypt_nodes`;

CREATE TABLE `brypt_nodes` (
	`id` varchar(32) NOT NULL,
	`serial_number` varchar(255) NOT NULL,
	`type` varchar(32),
	`created_on` datetime,
	`registered_on` datetime,
	`registered_to` varchar(32),
	`connected_network`varchar(32),
	--TODO: Add foreign key constraints once other tables created
	PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=latin1;
--MongoDB???
