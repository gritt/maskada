#
# CATEGORY
SET @FOOD = "Food";
SET @HEALTH = "Health";
SET @ENTERTAINMENT = "Entertainment";
SET @WORK = "Work";
SET @HOME = "Home";
#
INSERT INTO `category` (`name`)
VALUES (@FOOD),
       (@HEALTH),
       (@ENTERTAINMENT),
       (@WORK),
       (@HOME);

#
# TRANSACTION
SET @DEBIT = 1;
SET @CREDIT = 2;
SET @INCOME = 3;
#
INSERT INTO `transaction` (`amount`, `type`, `category`, `description`)
VALUES ("99", @CREDIT, @ENTERTAINMENT, NULL),
       ("11", @CREDIT, @FOOD, NULL),
       ("32", @CREDIT, @FOOD, NULL),
       ("5300", @INCOME, @WORK, NULL),
       ("129", @DEBIT, @HOME, "Internet"),
       ("129", @DEBIT, @HOME, "Electricity");