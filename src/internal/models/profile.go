package models

type Profile struct {
    ID int
    PictureID *int
    PhoneNumber *string
    Address *string
    Name *string
    FamilyName *string
    MiddleName *string
}
