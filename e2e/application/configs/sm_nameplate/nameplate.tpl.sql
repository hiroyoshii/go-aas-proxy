select
    digital_product.serial_number as SerialNumber,
    digital_product.country_of_origin as ProductCountryOfOrigin,
    digital_product.construction_year as YearOfConstruction
from
    digital_product
WHERE
    product_id = '{{ last (splitList "/" .AasID) }}'
;