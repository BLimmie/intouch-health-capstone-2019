import React from 'react'

import { Grid, Box, makeStyles, createStyles, Theme } from '@material-ui/core'
import Webcam from 'react-webcam'
import { Redirect, RouteComponentProps } from 'react-router-dom'
import { Room } from 'twilio-video'
import { RoomTextField } from '../Video/RoomTextField'
import { getRoom, setRemoteVideo } from '../Video/Twilio'
import VideoControls, { avStateInterface, defaultAVState } from '../Video/VideoControls'
import { WebcamWithControls } from '../Video/WebcamWithControls'

const useStyles = makeStyles((theme: Theme) =>
  createStyles({
    root: {
      background: '#16001f'
    },
    upperHalf: {
      height: "55vh"
    },
    their_video: {
      '& video' : {
        width: '100%',
        height: '90%',
      }
    },
    lowerHalf: {
      height: "45vh",
    }
  }),
);

type LinkParams = {link: string}
export default function PatientInterface({ match }: RouteComponentProps<LinkParams>)  {
  
  const classes = useStyles();
  const webcamRef:  React.RefObject<Webcam> = React.useRef(null)
  const [avState, setAvState] = React.useState<avStateInterface>(defaultAVState)
  
  const [roomField, setRoomField] = React.useState<string>(match.params.link)
  const [videoRoom, setVideoRoom] = React.useState<Room>()
  const [endChat, setEndChat] = React.useState<boolean>(false)

  // const joinRoom = () => {
  //   const localTrackStream = (webcamRef.current as Webcam ).stream.getTracks()
  //   getRoom( roomField, localTrackStream, "Patient")
  //     .then((room: Room) => {
  //       setVideoRoom(room)
  //       setRemoteVideo(room, endSession)
  //     })
  //     .catch(() => {setRoomField("Room name does not exist, try again")})
  // }
  const joinRoom = () => {
    if (!videoRoom) {
      if (navigator.mediaDevices.getUserMedia) { 
        navigator.mediaDevices.getUserMedia({video: true, audio: true}).then(stream => {
          const localTrackStream = stream.getTracks()
          getRoom(roomField, localTrackStream, "Patient")
            .then((room: Room) => {
              setVideoRoom(room)
              setRemoteVideo(room, endSession)
            })
            .catch(() => {setRoomField("Room name does not exist, try again")})
          }).catch(console.log)
      }
    }
  }
  joinRoom()

  const endSession: Function = () => {
    if(videoRoom)
      videoRoom.disconnect()
    setEndChat(true)
  }

  if(endChat)
    return <Redirect to='/client/' />
  
  return (
    <div className={classes.root} >
    <Grid container className={classes.upperHalf} >
      <Grid item xs={1} />
      {
        // videoRoom &&
        <Grid item className={classes.their_video}  id="remote" xs = {5} />
      } 
      <Grid item xs={1} />
      <Grid item xs={5} />
    </Grid>
    <Grid container className={classes.lowerHalf}>
      <Grid item xs={3} /*Beginning of the lower half*/ >
        <WebcamWithControls
          webcamRef={webcamRef}
          avState={avState}
          sendScreenshots={true}
          room={videoRoom}
        />
      </Grid>
      <Grid item xs={4}> 
        <VideoControls
          endSession={() => endSession()}
          avState={avState}
          setAVState={setAvState}
        />
      </Grid>
      <Grid item xs={5}>
        <Box
          className={classes.lowerHalf}
          border = {8}
          borderColor = "white"
          borderRadius = "0%"
        >
        </Box>
      </Grid>
    </Grid>
  </div>
  );

}