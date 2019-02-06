ALTER TABLE `delivery_order`
CHANGE `status` `status` enum('active','finish','cancelled') COLLATE 'latin1_swedish_ci' NOT NULL DEFAULT 'active' AFTER `number_seal`;