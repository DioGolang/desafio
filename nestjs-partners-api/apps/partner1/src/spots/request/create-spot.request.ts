import {IsNotEmpty, IsString, MaxLength} from "class-validator";

export class CreateSpotRequest {
  @IsNotEmpty({message: 'the name field is empty'})
  @IsString({message: 'invalid name'})
  @MaxLength(255, {message: 'maximum 255 characters'})
  name: string;
}
