select
    manufacturer.name as ManufacturerName,
    manufacturer.gl_number as GLNOfManufacturer,
    manufacturer.supplier_name as NameOfSupplier,
    product.serial_number as SerialNumber,
    product.batch_number as BatchNumber,
    contact_info.email as Email,
    contact_info.url as URL,
    contact_info.phone_number as PhoneNumber,
    contact_info.fax as Fax
from
    manufacturer
    left join product on product.manufacturer_id = manufacturer.id
    left join contact_info on contact_info.manufacturer_id = manufacturer.id
WHERE
    product.id = '{{ last (splitList "/" .AasID) }}'
;