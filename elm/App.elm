{- This app is the basic counter app. You can increment and decrement the count
like normal. The big difference is that the current count shows up in the URL.
Try changing the URL by hand. If you change it to a number, the app will go
there. If you change it to some invalid address, the app will recover in a
reasonable way.
-}

import Html exposing (..)
import Html.Attributes exposing (href)
import Navigation
import Query


main : Program Never
main =
  Navigation.program (Navigation.makeParser loginParser)
    { init = init
    , view = view
    , update = update
    , urlUpdate = urlUpdate
    , subscriptions = subscriptions
    }



-- URL PARSERS


loginParser : Navigation.Location -> List (String, Maybe String)
loginParser location =
  Query.parseQuery location.search


-- MODEL


type alias Model =
  { id : Maybe String
  }


init : List (String, Maybe String) -> (Model, Cmd Msg)
init result =
  urlUpdate result { id = Nothing }



-- UPDATE


type Msg = Login


{-| A relatively normal update function. The only notable thing here is that we
are commanding a new URL to be added to the browser history. This changes the
address bar and lets us use the browser&rsquo;s back button to go back to
previous pages.
-}
update : Msg -> Model -> (Model, Cmd Msg)
update msg model =
  ( model, Cmd.none )

{-| The URL is turned into a result. If the URL is valid, we just update our
model to the new count. If it is not a valid URL, we modify the URL to make
sense.
-}
urlUpdate : List (String, Maybe String) -> Model -> (Model, Cmd Msg)
urlUpdate list model =
  ( model, Cmd.none )


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
