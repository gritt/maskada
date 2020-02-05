DROP TABLE IF EXISTS `transaction`;
DROP TABLE IF EXISTS `category`;

CREATE TABLE `category`
(
    `name` VARCHAR(80) UNIQUE NOT NULL,
    PRIMARY KEY (`name`)
) ENGINE = InnoDB
  DEFAULT CHARSET = utf8mb4
  COLLATE = utf8mb4_unicode_ci;

CREATE TABLE `transaction`
(
    `id`          INTEGER(11) NOT NULL AUTO_INCREMENT,
    `amount`      INTEGER(11) NOT NULL DEFAULT 0,
    `type`        INTEGER(11) NOT NULL,
    `category`    VARCHAR(80) NOT NULL,
    CONSTRAINT `fk_category`
        FOREIGN KEY (`category`) REFERENCES `category` (`name`)
            ON DELETE RESTRICT
            ON UPDATE CASCADE,
    `description` VARCHAR(80) NULL,
    `date`        TIMESTAMP   NOT NULL DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (`id`)
) ENGINE = InnoDB
  DEFAULT CHARSET = utf8mb4
  COLLATE = utf8mb4_unicode_ci;
