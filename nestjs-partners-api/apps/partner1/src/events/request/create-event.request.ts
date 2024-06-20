import {IsInt, IsISO8601, IsNotEmpty, IsString, MaxLength, Min} from "class-validator";

export class CreateEventRequest {
  @IsNotEmpty({message: 'the name field is empty'})
  @IsString({message: 'invalid name'})
  @MaxLength(255, {message: 'maximum 255 characters'})
  name: string;

  @IsNotEmpty({message: 'the description field is empty'})
  @IsString({message: 'invalid description'})
  @MaxLength(255, {message: 'maximum 255 characters'})
  description: string;

  @IsNotEmpty({message: 'the date field is empty'})
  @IsISO8601()
  @IsString({message: 'invalid field'})
  date: string;

  @IsNotEmpty({message: 'the price field is empty'})
  @IsInt({message:'number invalided'})
  @Min(0, {message:'numbers must be positive'})
  price: number;
}
