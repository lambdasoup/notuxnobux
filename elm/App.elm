{- This app is the basic counter app. You can increment and decrement the count
   like normal. The big difference is that the current count shows up in the URL.
   Try changing the URL by hand. If you change it to a number, the app will go
   there. If you change it to some invalid address, the app will recover in a
   reasonable way.
-}


module Main exposing (..)

import Html exposing (..)
import Html.Attributes exposing (href)
import Navigation
import String
import Task
import Http
import Json.Decode as Json


main : Program Never
main =
    Navigation.program (Navigation.makeParser queryParser)
        { init = init
        , view = view
        , update = update
        , urlUpdate = urlUpdate
        , subscriptions = subscriptions
        }



-- URL PARSERS


queryParser : Navigation.Location -> String
queryParser location =
    location.search



-- MODEL


type alias Model =
    { id : Maybe String
    }


init : String -> ( Model, Cmd Msg )
init result =
    urlUpdate result { id = Nothing }



-- UPDATE


type Msg
    = GetJwtSucceed String
    | GetJwtFail Http.Error


{-| A relatively normal update function. The only notable thing here is that we
are commanding a new URL to be added to the browser history. This changes the
address bar and lets us use the browser&rsquo;s back button to go back to
previous pages.
-}
update : Msg -> Model -> ( Model, Cmd Msg )
update msg model =
    ( model, Cmd.none )


{-| The URL is turned into a result. If the URL is valid, we just update our
model to the new count. If it is not a valid URL, we modify the URL to make
sense.
-}
urlUpdate : String -> Model -> ( Model, Cmd Msg )
urlUpdate query model =
    if String.isEmpty query then
        ( { id = Nothing }, Cmd.none )
    else
        ( { id = Just "pending" }, login query )


login : String -> Cmd Msg
login query =
    let
        url =
            "jwt" ++ query
    in
        Task.perform GetJwtFail GetJwtSucceed (Http.get nop url)


nop : Json.Decoder String
nop =
    Json.at [ "data", "image_url" ] Json.string



-- SUBSCRIPTIONS


subscriptions : Model -> Sub Msg
subscriptions model =
    Sub.none



-- VIEW


view : Model -> Html Msg
view model =
    div []
        [ div [] [ text (toString model) ]
        , a [ href "login" ] [ text "Login" ]
        ]
